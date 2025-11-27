import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ShoppingCartService } from '../service/shopping-cart.service';
import { AuthService } from '../../infrastructure/auth.service';
import { CartResponse, CheckoutResponse } from '../model/model';
import { RouterLink } from '@angular/router';


@Component({
  selector: 'app-shopping-cart',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './shopping-cart.component.html',
  styleUrls: ['./shopping-cart.component.css']
})
export class ShoppingCartComponent implements OnInit {
  cart: CartResponse | null = null;
  isLoading = false;
  errorMsg = '';

  // To store checkout result
  purchaseSuccess = false;
  purchasedTokens: any[] = [];

  userId: string = '';

  constructor(
    private cartService: ShoppingCartService,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    // 1. Extract User ID from Token
    this.userId = this.getUserIdFromToken();

    if (this.userId) {
      this.loadCart();
    } else {
      this.errorMsg = 'Please log in to view your cart.';
    }
  }

  getUserIdFromToken(): string {
    const token = this.authService.getToken();
    if (!token) return '';
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.id || payload.sub || '';
    } catch (e) {
      return '';
    }
  }

  loadCart(): void {
    this.isLoading = true;
    this.cartService.getCart(this.userId).subscribe({
      next: (data) => {
        this.cart = data;
        this.isLoading = false;
      },
      error: (err) => {
        console.error(err);
        this.errorMsg = 'Failed to load cart.';
        this.isLoading = false;
      }
    });
  }

  removeItem(tourId: string): void {
    if (!confirm('Are you sure you want to remove this item?')) return;

    this.cartService.removeItem(this.userId, tourId).subscribe({
      next: (updatedCart) => {
        this.cart = updatedCart;
      },
      error: (err) => {
        alert('Failed to remove item');
      }
    });
  }

  checkout(): void {
    if (!this.cart || this.cart.items.length === 0) return;
    this.isLoading = true;

    this.cartService.checkout(this.userId).subscribe({
      next: (response: CheckoutResponse) => {
        this.purchaseSuccess = true;
        this.purchasedTokens = response.tokens;
        this.cart = null; // Clear cart view
        this.isLoading = false;
        alert('purchase success.');
      },
      error: (err) => {
        this.errorMsg = 'Checkout failed. Please try again.';
        this.isLoading = false;
      }
    });
  }
}