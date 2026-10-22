package mapper

import (
	"time"

	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/entity"
)

// PostSaveWithProductToListing maps the listing subset of PostSaveWithProduct (fields 2–9).
func PostSaveWithProductToListing(post *bdspropb.PostSaveWithProduct, ownerProfileID uint64) *entity.Listing {
	if post == nil {
		return nil
	}
	numDate := post.GetNumDate()
	if numDate <= 0 {
		numDate = 15
	}
	if numDate > 30 {
		numDate = 30
	}

	listing := &entity.Listing{
		OwnerID:         ownerProfileID,
		Title:           post.GetTitle(),
		Content:         post.GetContent(),
		PostPrice:       post.GetPostPrice(),
		TransactionType: post.GetTransactionType(),
		Visibility:      post.GetVisibility(),
		NumDate:         numDate,
	}
	if post.GetProductId() != 0 {
		listing.ProductID = post.GetProductId()
	}

	if exp := post.GetExpiredAt(); exp > 0 {
		// proto int32: unix seconds (millis do not fit int32).
		t := time.Unix(int64(exp), 0)
		listing.ExpiredAt = &t
	} else {
		t := time.Now().Add(time.Duration(numDate) * 24 * time.Hour)
		listing.ExpiredAt = &t
	}
	return listing
}

// ListingToPost maps a listing row to bdspropb.Post (fields supported by v4 listings table).
func ListingToPost(model *entity.Listing) *bdspropb.Post {
	if model == nil {
		return nil
	}
	out := &bdspropb.Post{
		Id:              model.ID,
		ProductId:       model.ProductID,
		Visibility:      uint32(model.Visibility),
		Content:         model.Content,
		TransactionType: uint32(model.TransactionType),
		Title:           model.Title,
		NumDate:         model.NumDate,
		PostPrice:       model.PostPrice,
		OwnerId:         model.OwnerID,
	}
	if model.ExpiredAt != nil {
		out.ExpiredAt = _utils.FormatTimeToString(model.ExpiredAt)
	}
	return out
}

// ListingsToPosts maps listing rows for list APIs.
func ListingsToPosts(rows []entity.Listing) []*bdspropb.Post {
	out := make([]*bdspropb.Post, 0, len(rows))
	for i := range rows {
		out = append(out, ListingToPost(&rows[i]))
	}
	return out
}
