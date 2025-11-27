import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { AuthService } from '../../infrastructure/auth.service';
import { CartResponse, CheckoutResponse } from '../model/model';

@Injectable({
  providedIn: 'root'
})
export class ShoppingCartService {
  // Base URL: http://localhost:8080/api/shopping-cart
  private apiUrl = 'http://localhost:8080/api/shopping-cart';

  private cartCountSubject = new BehaviorSubject<number>(0);
  public cartCount$ = this.cartCountSubject.asObservable();

  constructor(
    private http: HttpClient,
    private authService: AuthService
  ) { }

  // Helper for Auth Headers
  private getAuthHeaders(): HttpHeaders {
    const token = this.authService.getToken();
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }

  private fetchCartCount() {
    const userId = this.authService.getMyId();
    if (!userId) return;

    this.http.get<CartResponse>(`${this.apiUrl}/${userId}`, { headers: this.getAuthHeaders() })
      .subscribe({
        next: (cart) => this.cartCountSubject.next(cart.items.length),
        error: () => this.cartCountSubject.next(0)
      });
  }


  // 1. GET CART
  getCart(userId: string): Observable<CartResponse> {
    return this.http.get<CartResponse>(`${this.apiUrl}/${userId}`, {
      headers: this.getAuthHeaders()
    }).pipe(
      tap(cart => this.cartCountSubject.next(cart.items.length))
    );
  }

  // 2. ADD ITEM (To be called from Tour List)
  addItem(userId: string, tourId?: string): Observable<CartResponse> {
    const body = {
      user_id: userId,
      tour_id: tourId
    };
    return this.http.post<CartResponse>(`${this.apiUrl}/items`, body, {
      headers: this.getAuthHeaders()
    }).pipe(
      // Update count after adding item
      tap(cart => this.cartCountSubject.next(cart.items.length))
    );
  }

  // 3. REMOVE ITEM
  removeItem(userId: string, tourId: string): Observable<CartResponse> {
    return this.http.delete<CartResponse>(`${this.apiUrl}/items/${userId}/${tourId}`, {
      headers: this.getAuthHeaders()
    }).pipe(
      // Update count after adding item
      tap(cart => this.cartCountSubject.next(cart.items.length))
    );
  }

  // 4. CHECKOUT
  checkout(userId: string): Observable<CheckoutResponse> {
    const body = { userId };
    return this.http.post<CheckoutResponse>(`${this.apiUrl}/checkout`, body, {
      headers: this.getAuthHeaders()
    }).pipe(
      // Update count after adding item
      tap(cart => this.cartCountSubject.next(0))
    );
  }
}