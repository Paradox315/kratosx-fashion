package repo

import (
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/google/wire"
)

const HTTP = "http://"

var codec = encoding.GetCodec("json")

var ProviderSet = wire.NewSet(
	NewClothesRepo,
	NewRecommendRepo,
	NewFeedbackRepo,
)
