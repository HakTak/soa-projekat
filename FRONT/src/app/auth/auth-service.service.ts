import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
// Importujemo modele
import { User, UserListResponse } from '../models/user.model';
// Importujemo infrastrukturni servis da bismo dobili token
import { AuthService } from '../infrastructure/auth.service';

@Injectable({
  providedIn: 'root'
})
export class AuthFeatureService {
  // Napomena: Klasu sam nazvao AuthFeatureService da se ne mesa sa AuthService iz infrastructure.
  // Ime fajla ostaje auth-service.service.ts

  private apiUrl = 'http://localhost:8080/api';

  // Ubrizgavamo AuthService da bismo koristili njegovu getToken() metodu
  constructor(private http: HttpClient, private authService: AuthService) { }

  // --- METODE PREBACENE OVDE ---

  getUsers(): Observable<UserListResponse> {
    const headers = this.getAuthHeaders();
    return this.http.get<UserListResponse>(`${this.apiUrl}/admin/users`, { headers });
  }

  blockUser(userId: string): Observable<User> {
    const headers = this.getAuthHeaders();
    return this.http.put<User>(`${this.apiUrl}/admin/block/${userId}`, {}, { headers });
  }

  // Helper za headere (koristi token iz Infra servisa)
  getAuthHeaders(): HttpHeaders {
    const token = this.authService.getToken();
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }
}