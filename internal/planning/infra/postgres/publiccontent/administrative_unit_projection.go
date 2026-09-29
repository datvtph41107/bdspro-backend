package publiccontent

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"tqd/internal/domain/publiccontent/model"

	"gorm.io/gorm"
)

var administrativeCodeSuffixPattern = regexp.MustCompile(`(?:^|-)([0-9A-Za-z]+)$`)

type administrativeUnitRow struct {
	EntityID       string
	SourceID       string
	UnitType       string
	FullName       string
	ShortName      string
	Code           string
	ParentID       string
	ParentFullName string
	ParentCode     string
	Latitude       float64
	Longitude      float64
	UpdatedAt      *time.Time
}

func normalizeAdministrativeIdentity(value string) (unitType string, identity string) {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "/")
	if index := strings.IndexAny(value, "?#"); index >= 0 {
		value = value[:index]
	}
	if index := strings.LastIndex(value, "/"); index >= 0 {
		value = value[index+1:]
	}
	value = strings.TrimSpace(value)
	if prefix, rest, ok := strings.Cut(value, ":"); ok {
		switch strings.ToLower(strings.TrimSpace(prefix)) {
		case "ward", "province":
			return strings.ToLower(strings.TrimSpace(prefix)), strings.TrimSpace(rest)
		}
	}
	return "", value
}

func administrativeCodeCandidate(identity string) string {
	identity = strings.TrimSpace(identity)
	match := administrativeCodeSuffixPattern.FindStringSubmatch(identity)
	if len(match) != 2 {
		return identity
	}
	return strings.TrimSpace(match[1])
}

func administrativeUnitSlug(name string, code string) (string, error) {
	nameSlug := slugify(name)
	codeSlug := slugify(code)
	if nameSlug == "" {
		return "", errors.New("administrative unit name cannot produce a canonical slug")
	}
	if codeSlug == "" {
		return "", errors.New("administrative unit code is required for a stable canonical slug")
	}
	return nameSlug + "-" + codeSlug, nil
}

func administrativeUnitTypeName(unitType string, fullName string) string {
	name := strings.TrimSpace(fullName)
	lowerName := strings.ToLower(name)
	for prefix, label := range map[string]string{
		"thành phố": "Thành phố",
		"tỉnh":      "Tỉnh",
		"phường":    "Phường",
		"xã":        "Xã",
		"thị trấn":  "Thị trấn",
	} {
		if strings.HasPrefix(lowerName, prefix+" ") || lowerName == prefix {
			return label
		}
	}
	if unitType == "province" {
		return "Tỉnh/Thành phố"
	}
	return "Phường/Xã"
}

func (r *Repository) queryAdministrativeUnits(
	ctx context.Context,
	unitType string,
	identity string,
) ([]administrativeUnitRow, error) {
	codeCandidate := administrativeCodeCandidate(identity)
	rows := make([]administrativeUnitRow, 0, 2)

	if unitType == "" || unitType == "ward" {
		var wardRows []administrativeUnitRow
		if err := r.db.WithContext(ctx).Raw(`
			SELECT CONCAT('ward:', w.id) AS entity_id,
			       w.id::text AS source_id,
			       'ward' AS unit_type,
			       COALESCE(w.full_name, '') AS full_name,
			       COALESCE(w.short_name, '') AS short_name,
			       COALESCE(w.code::text, '') AS code,
			       COALESCE(p.id::text, '') AS parent_id,
			       COALESCE(p.full_name, '') AS parent_full_name,
			       COALESCE(p.code::text, '') AS parent_code,
			       COALESCE(w.lat, 0)::float8 AS latitude,
			       COALESCE(w.lng, 0)::float8 AS longitude,
			       NULL::timestamptz AS updated_at
			FROM ward_v2 w
			LEFT JOIN province_v2 p ON p.id = w.province_id
			WHERE w.id::text = ?
			   OR LOWER(COALESCE(w.code::text, '')) = LOWER(?)
			   OR LOWER(COALESCE(w.full_name, '')) = LOWER(?)
			   OR LOWER(COALESCE(w.short_name, '')) = LOWER(?)
			ORDER BY w.created_at DESC NULLS LAST, w.id::text ASC
			LIMIT 3`, identity, codeCandidate, identity, identity).Scan(&wardRows).Error; err != nil {
			return nil, fmt.Errorf("query ward projection identity %q: %w", identity, err)
		}
		rows = append(rows, wardRows...)
	}

	if unitType == "" || unitType == "province" {
		var provinceRows []administrativeUnitRow
		if err := r.db.WithContext(ctx).Raw(`
			SELECT CONCAT('province:', p.id) AS entity_id,
			       p.id::text AS source_id,
			       'province' AS unit_type,
			       COALESCE(p.full_name, '') AS full_name,
			       COALESCE(p.short_name, '') AS short_name,
			       COALESCE(p.code::text, '') AS code,
			       '' AS parent_id,
			       '' AS parent_full_name,
			       '' AS parent_code,
			       COALESCE(p.lat, 0)::float8 AS latitude,
			       COALESCE(p.lng, 0)::float8 AS longitude,
			       p.updated_at
			FROM province_v2 p
			WHERE p.id::text = ?
			   OR LOWER(COALESCE(p.code::text, '')) = LOWER(?)
			   OR LOWER(COALESCE(p.full_name, '')) = LOWER(?)
			   OR LOWER(COALESCE(p.short_name, '')) = LOWER(?)
			ORDER BY p.updated_at DESC NULLS LAST
			LIMIT 3`, identity, codeCandidate, identity, identity).Scan(&provinceRows).Error; err != nil {
			return nil, fmt.Errorf("query province projection identity %q: %w", identity, err)
		}
		rows = append(rows, provinceRows...)
	}
	return rows, nil
}

func selectAdministrativeUnitRow(
	rows []administrativeUnitRow,
	unitType string,
	identity string,
) (*administrativeUnitRow, error) {
	identity = strings.TrimSpace(identity)
	identityLower := strings.ToLower(identity)
	matches := make([]administrativeUnitRow, 0, len(rows))
	for _, row := range rows {
		if unitType != "" && row.UnitType != unitType {
			continue
		}
		slug, err := administrativeUnitSlug(row.FullName, row.Code)
		if err != nil {
			continue
		}
		for _, candidate := range []string{
			row.EntityID,
			row.SourceID,
			row.Code,
			row.FullName,
			row.ShortName,
			slug,
		} {
			if strings.ToLower(strings.TrimSpace(candidate)) == identityLower {
				matches = append(matches, row)
				break
			}
		}
	}
	if len(matches) == 1 {
		return &matches[0], nil
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf(
			"administrative unit identity %q is ambiguous; use ward:<id> or province:<id>",
			identity,
		)
	}
	if len(rows) == 1 {
		return &rows[0], nil
	}
	if len(rows) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return nil, fmt.Errorf(
		"administrative unit identity %q is ambiguous; use a canonical slug or typed identity",
		identity,
	)
}

type administrativePlanningCountRow struct {
	JurisdictionMatches  int64
	PlanningProjectCount int64
}

type administrativePlanningProjectRow struct {
	ID           uint64
	Slug         string
	Code         string
	Name         string
	Summary      string
	ProjectCount int64
}

func (r *Repository) administrativePlanningProjectCount(
	ctx context.Context,
	fullName string,
	shortName string,
) (*int64, error) {
	fullName = strings.TrimSpace(fullName)
	shortName = strings.TrimSpace(shortName)
	if fullName == "" && shortName == "" {
		return nil, nil
	}

	var row administrativePlanningCountRow
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT j.id) AS jurisdiction_matches,
		       COUNT(DISTINCT pp.id) AS planning_project_count
		FROM qh_jurisdictions j
		LEFT JOIN qh_planning_projects pp
		  ON pp.jurisdiction_id = j.id
		 AND pp.deleted_at IS NULL
		WHERE j.deleted_at IS NULL
		  AND (
		    LOWER(TRIM(j.name)) = LOWER(TRIM(?))
		    OR (? <> '' AND LOWER(TRIM(j.name)) = LOWER(TRIM(?)))
		  )`, fullName, shortName, shortName).Scan(&row).Error; err != nil {
		return nil, err
	}
	if row.JurisdictionMatches != 1 {
		return nil, nil
	}
	return &row.PlanningProjectCount, nil
}

func (r *Repository) administrativePlanningProjects(
	ctx context.Context,
	fullName string,
	shortName string,
) ([]domain.AdministrativePlanningProject, *int64, error) {
	fullName = strings.TrimSpace(fullName)
	shortName = strings.TrimSpace(shortName)
	if fullName == "" && shortName == "" {
		return nil, nil, nil
	}

	rows := make([]administrativePlanningProjectRow, 0, 12)
	if err := r.db.WithContext(ctx).Raw(`
		WITH matched_jurisdiction AS (
			SELECT j.id
			FROM qh_jurisdictions j
			WHERE j.deleted_at IS NULL
			  AND (
			    LOWER(TRIM(j.name)) = LOWER(TRIM(?))
			    OR (? <> '' AND LOWER(TRIM(j.name)) = LOWER(TRIM(?)))
			  )
		), unique_jurisdiction AS (
			SELECT MIN(id) AS id
			FROM matched_jurisdiction
			HAVING COUNT(*) = 1
		), projects AS (
			SELECT pp.id,
			       COALESCE(pp.slug, '') AS slug,
			       COALESCE(pp.code, '') AS code,
			       COALESCE(pp.name, '') AS name,
			       COALESCE(pp.summary, '') AS summary,
			       COUNT(*) OVER() AS project_count
			FROM qh_planning_projects pp
			JOIN unique_jurisdiction uj ON uj.id = pp.jurisdiction_id
			WHERE pp.deleted_at IS NULL
			ORDER BY pp.updated_at DESC NULLS LAST, pp.id DESC
			LIMIT 12
		)
		SELECT * FROM projects`,
		fullName,
		shortName,
		shortName,
	).Scan(&rows).Error; err != nil {
		return nil, nil, err
	}

	projects := make([]domain.AdministrativePlanningProject, 0, len(rows))
	var total *int64
	for _, row := range rows {
		if total == nil {
			count := row.ProjectCount
			total = &count
		}

		slug := strings.TrimSpace(row.Slug)
		name := strings.TrimSpace(row.Name)
		if row.ID == 0 || slug == "" || name == "" {
			continue
		}

		entityID := fmt.Sprintf("planning_project:%d", row.ID)
		projects = append(projects, domain.AdministrativePlanningProject{
			EntityID:   entityID,
			Slug:       slug,
			Code:       strings.TrimSpace(row.Code),
			Name:       name,
			Summary:    strings.TrimSpace(row.Summary),
			PublicPath: "/do-an-quy-hoach/" + slug,
			MapPath:    "/ban-do?focus=" + url.QueryEscape(entityID),
		})
	}

	if total == nil {
		zero := int64(0)
		total = &zero
	}

	return projects, total, nil
}

func (r *Repository) GetAdministrativeUnitProjection(
	ctx context.Context,
	identity string,
) (*domain.AdministrativeUnitProjection, error) {
	unitType, normalizedIdentity := normalizeAdministrativeIdentity(identity)
	if normalizedIdentity == "" {
		return nil, errors.New("administrative unit identity is required")
	}

	rows, err := r.queryAdministrativeUnits(
		ctx,
		unitType,
		normalizedIdentity,
	)
	if err != nil {
		return nil, err
	}

	row, err := selectAdministrativeUnitRow(
		rows,
		unitType,
		normalizedIdentity,
	)
	if err != nil {
		return nil, err
	}

	slug, err := administrativeUnitSlug(row.FullName, row.Code)
	if err != nil {
		return nil, fmt.Errorf(
			"build administrative unit canonical slug: %w",
			err,
		)
	}

	canonicalPath := "/dia-ban/" + slug
	canonicalKey := "administrative_unit:" +
		row.UnitType +
		":" +
		strings.TrimSpace(row.Code)

	limitations := []string{
		"Nguồn hiện tại chỉ có điểm trung tâm hành chính; chưa có polygon địa giới để xác nhận point-in-boundary.",
		"Thông tin quy hoạch cần được đối chiếu với lớp bản đồ và hồ sơ pháp lý gốc.",
	}

	if row.UpdatedAt == nil {
		limitations = append(
			limitations,
			"Nguồn ward_v2 hiện chưa có updated_at; source.updatedAt được để trống thay vì suy diễn từ created_at.",
		)
	}

	dataQuality := "partial"

	geometry := preview(
		"Point",
		row.Latitude,
		row.Longitude,
		0,
		0,
		0,
		0,
	)
	if geometry.Centroid != nil {
		geometry.Completeness = "centroid_only"
		geometry.FitMode = "center"
	}

	var location *domain.Point
	if row.Latitude != 0 || row.Longitude != 0 {
		location = &domain.Point{
			Latitude:  row.Latitude,
			Longitude: row.Longitude,
		}
	} else {
		limitations = append(
			limitations,
			"Đơn vị hành chính chưa có tọa độ trung tâm trong nguồn hiện tại.",
		)
	}

	relatedProjects, projectCount, countErr :=
		r.administrativePlanningProjects(
			ctx,
			row.FullName,
			row.ShortName,
		)

	availability := "name_matched"
	if countErr != nil || projectCount == nil {
		projectCount = nil
		relatedProjects = nil
		availability = "unavailable"

		limitations = append(
			limitations,
			"Chưa có mapping administrative code ↔ planning jurisdiction đủ tin cậy để công bố số lượng đồ án liên kết.",
		)
	} else {
		limitations = append(
			limitations,
			"Số lượng đồ án hiện được liên kết bằng tên jurisdiction khớp chính xác; cần thay bằng mapping administrative code trước production rollout.",
		)
	}

	hierarchy := make(
		[]domain.AdministrativeUnitHierarchyItem,
		0,
		2,
	)

	var parent *domain.AdministrativeUnitParent
	if row.UnitType == "ward" &&
		strings.TrimSpace(row.ParentFullName) != "" {
		parentSlug, parentSlugErr :=
			administrativeUnitSlug(
				row.ParentFullName,
				row.ParentCode,
			)

		if parentSlugErr == nil {
			parentPath := "/dia-ban/" + parentSlug

			parent = &domain.AdministrativeUnitParent{
				EntityID:           "province:" + row.ParentID,
				CanonicalKey:       "administrative_unit:province:" + row.ParentCode,
				Slug:               parentSlug,
				CanonicalPath:      parentPath,
				UnitType:           "province",
				Name:               row.ParentFullName,
				AdministrativeCode: row.ParentCode,
			}

			hierarchy = append(
				hierarchy,
				domain.AdministrativeUnitHierarchyItem{
					EntityID:           parent.EntityID,
					UnitType:           parent.UnitType,
					Name:               parent.Name,
					AdministrativeCode: parent.AdministrativeCode,
					CanonicalPath:      parent.CanonicalPath,
				},
			)
		} else {
			limitations = append(
				limitations,
				"Đơn vị cha chưa có administrative code ổn định để tạo canonical path.",
			)
		}
	}

	hierarchy = append(
		hierarchy,
		domain.AdministrativeUnitHierarchyItem{
			EntityID:           row.EntityID,
			UnitType:           row.UnitType,
			Name:               row.FullName,
			AdministrativeCode: row.Code,
			CanonicalPath:      canonicalPath,
		},
	)

	planningSummary :=
		"Tra cứu bối cảnh quy hoạch và các lớp dữ liệu liên quan của " +
			row.FullName +
			" trên bản đồ QHPro."

	return &domain.AdministrativeUnitProjection{
		SchemaVersion:      domain.AdministrativeUnitProjectionSchemaVersion,
		EntityID:           row.EntityID,
		CanonicalKey:       canonicalKey,
		Slug:               slug,
		CanonicalPath:      canonicalPath,
		UnitType:           row.UnitType,
		UnitTypeName:       administrativeUnitTypeName(row.UnitType, row.FullName),
		Name:               row.FullName,
		DisplayName:        row.FullName,
		AdministrativeCode: row.Code,
		Parent:             parent,
		Hierarchy:          hierarchy,
		Location:           location,
		GeometrySummary:    geometry,
		PlanningContext: domain.AdministrativePlanningContext{
			Summary:              planningSummary,
			PlanningProjectCount: projectCount,
			DataAvailability:     availability,
			RelatedProjects:      relatedProjects,
		},
		Source: domain.ProjectionSource{
			System: "tqd-service",
			Dataset: map[string]string{
				"province": "province_v2",
				"ward":     "ward_v2",
			}[row.UnitType],
			Authority:      "aggregated",
			AuthorityName:  "Dữ liệu địa giới hành chính được đồng bộ vào TQD",
			AuthorityClass: "government_derived",
			DataQuality:    dataQuality,
			UpdatedAt:      row.UpdatedAt,
			Limitations:    limitations,
		},
	}, nil
}
