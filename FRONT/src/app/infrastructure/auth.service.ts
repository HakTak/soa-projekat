import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable, tap } from 'rxjs';

export interface LoginRequest { username: string; password: string; }
export interface AuthResponse { token: string; }
export interface RegisterRequest { username: string; password: string; email: string; role: string; }

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:8080/api/auth';

  // 1. Inicijalizacija (Privatno stanje)
  private userState = new BehaviorSubject<boolean>(!!localStorage.getItem('jwt'));

  // 2. JAVNI Observable na koji se Navbar kaci
  // Dodao sam 'public' da budemo sigurni da ga druge komponente vide
  public userState$: Observable<boolean> = this.userState.asObservable();

  constructor(private http: HttpClient) { }

  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.apiUrl}/login`, credentials)
      .pipe(
        tap(response => {
          if (response.token) {
            localStorage.setItem('jwt', response.token);
            this.userState.next(true); // Javljamo da je ulogovan
          }
        })
      );
  }

  register(data: RegisterRequest): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, data);
  }

  logout() {
    localStorage.removeItem('jwt');
    this.userState.next(false); // Javljamo da je izlogovan
  }

  getToken(): string | null {
    return localStorage.getItem('jwt');
  }

  isLoggedIn(): boolean {
    return this.userState.value;
  }
}