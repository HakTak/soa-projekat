export interface CartItem {
    tourId: string;
    tourName: string;
    price: number;
}

export interface CartResponse {
    userId: string;
    items: CartItem[];
    totalPrice: number;
}

export interface PurchaseToken {
    tokenId: string;
    tourId: string;
    userId: string;
}

export interface CheckoutResponse {
    success: boolean;
    tokens: PurchaseToken[];
}