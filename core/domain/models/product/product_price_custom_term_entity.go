package product

type ProductPriceCustomTermEntity struct {
	price ProductPriceValue
	term  int // 契約更新サイクル日数
}
