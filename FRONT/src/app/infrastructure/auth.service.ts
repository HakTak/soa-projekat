import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';

// Interfejsi za Auth ostaju ovde
export interface LoginRequest { username: string; password: string; }
export interface AuthResponse { token: string; }
export interface RegisterRequest { username: string; password: string; email: string; role: string; }

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:8080/api/auth';

  private userState = new BehaviorSubject<boolean>(!!localStorage.getItem('jwt'));
  public userState$ = this.userState.asObservable();

  constructor(private http: HttpClient) { }

  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.apiUrl}/login`, credentials)
      .pipe(
        tap(response => {
          if (response.token) {
            localStorage.setItem('jwt', response.token);
            this.userState.next(true);
          }
        })
      );
  }

  register(data: RegisterRequest): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, data);
  }

  logout() {
    localStorage.removeItem('jwt');
    this.userState.next(false);
  }

  // --- HELPERI KOJI TREBAJU DRUGIM SERVISIMA ---

  getToken(): string | null {
    return localStorage.getItem('jwt');
  }

  // // Ovu proveru ostavljamo ovde jer zavisi od dekodiranja tokena
  // isAdmin(): boolean {
  //   const token = this.getToken();
  //   if (!token) return false;
  //   try {
  //     const payload = JSON.parse(atob(token.split('.')[1]));
  //     return payload.role === 'ADMIN';
  //   } catch (e) {
  //     return false;
  //   }
  // }

  hasRole(): string {
    const token = this.getToken();
    if (!token) return "UNAUTHORIZED";
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.role;
    } catch (e) {
      return "ROLE_ERROR"
    }
  }

  isLoggedIn(): boolean {
    return this.userState.value;
  }

  // Helper za headere (koristi token iz Infra servisa)
  public getAuthHeaders(): HttpHeaders {
    const token = this.getToken();
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }

  public getMyId(): string {
    const token = this.getToken()
    if (!token) return "Unknowen user token";
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      console.log("EVO GA ID MOJ: " + payload.id)
      return payload.id;
    } catch (e) {
      return "Error user token";
    }
  }

  public getMyUsername(): string {
    const token = this.getToken()
    if (!token) return "Unknowen user token"; 
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      console.log("EVO GA USERNAME MOJ: " + payload.username)
      return payload.username;
    } catch (e) {
      return "Error user token";
    }
  }
}