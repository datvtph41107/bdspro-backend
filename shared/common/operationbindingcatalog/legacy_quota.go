package operationbindingcatalog

import (
	_operationbinding "common/operationbinding"
)

/** LegacyQuotaReviewV1 records every legacy quota transport decision. */
func LegacyQuotaReviewV1() _operationbinding.Catalog {
	return _operationbinding.Catalog{
		Version: "1.0.0",
		Bindings: []_operationbinding.Binding{
			{
				Code:        "legacy.tqd_create_batch_parcels",
				ModuleOwner: "parcel",
				Transports: []_operationbinding.Transport{
					_operationbinding.GRPCUnary("/tqdpb.ParcelService/CreateBatchParcels"),
					_operationbinding.HTTPRoute("POST", "/v2/tqd/parcels/batch"),
				},
				Disposition: _operationbinding.DispositionReviewRequired,
				Description: "Batch quantity, caller class, and business success must be approved before operation mapping.",
			},
			{
				Code:        "legacy.tqd_find_parcel_by_location",
				ModuleOwner: "parcel",
				Transports: []_operationbinding.Transport{
					_operationbinding.GRPCUnary("/tqdpb.ParcelService/FindParcelByLocation"),
					_operationbinding.HTTPRoute("POST", "/v2/tqd/parcels/by-location"),
				},
				Disposition: _operationbinding.DispositionReviewRequired,
				Description: "The frontend marks this parcel-only endpoint as legacy fallback, so it cannot define the current business operation by itself.",
			},
			{
				Code:        "legacy.tqd_find_parcels_by_polygon",
				ModuleOwner: "parcel",
				Transports: []_operationbinding.Transport{
					_operationbinding.GRPCUnary("/tqdpb.ParcelService/FindParcelsByPolygon"),
					_operationbinding.HTTPRoute("POST", "/v2/tqd/parcels/by-polygon"),
				},
				Disposition: _operationbinding.DispositionReviewRequired,
				Description: "Zoom, intersection, and analysis flags change the response semantics, so a domain resolver or separate operations require approval.",
			},
			{
				Code:        "legacy.tqd_checking_polygon",
				ModuleOwner: "planning",
				Transports: []_operationbinding.Transport{
					_operationbinding.GRPCUnary("/tqdpb.ParcelService/CheckingPolygon"),
					_operationbinding.HTTPRoute("POST", "/v2/tqd/public/qh/checking-polygon"),
				},
				Disposition: _operationbinding.DispositionReviewRequired,
				Description: "The public draw workflow may be discovery, validation, or analysis and needs an approved customer-value boundary.",
			},
			{
				Code:        "legacy.tqd_get_parcel_planning",
				ModuleOwner: "planning",
				Transports: []_operationbinding.Transport{
					_operationbinding.GRPCUnary("/tqdpb.ParcelService/GetParcelPlanning"),
					_operationbinding.HTTPRoute("GET", "/v2/tqd/parcels/{parcelId}/planning"),
				},
				Disposition: _operationbinding.DispositionReviewRequired,
				Description: "The protobuf marks this RPC legacy while the frontend uses view and mode variants, so no single direct mapping is defensible yet.",
			},
			{
				Code:        "legacy.user_get_profile",
				ModuleOwner: "profile",
				Transports: []_operationbinding.Transport{
					_operationbinding.GRPCUnary("/userpb.ProfileService/GetProfileInfo"),
					_operationbinding.HTTPRoute("GET", "/v2/user/profile/info/{id}"),
					_operationbinding.GRPCUnary("/userpb.ProfileService/GetProfileInfoV3"),
					_operationbinding.HTTPRoute("GET", "/v3/user/profile/{access}/{id}"),
					_operationbinding.GRPCUnary("/userpb.ProfileService/GetProfileInfoMe"),
					_operationbinding.HTTPRoute("GET", "/v2/user/profile/info/me"),
				},
				Disposition: _operationbinding.DispositionReviewRequired,
				Description: "The legacy policy key is ambiguous across public, versioned, and current-user profile reads.",
			},
			{
				Code:        "legacy.user_upload_avatar",
				ModuleOwner: "profile",
				Disposition: _operationbinding.DispositionSourceMissing,
				Description: "No UploadAvatar RPC exists in the current protobuf schema.",
			},
		},
	}
}
