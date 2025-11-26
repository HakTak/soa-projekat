import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Profile, UpdateProfileRequest, ProfileStats } from '../model/profile';
import { AuthService } from '../../auth.service';

@Injectable({
  providedIn: 'root'
})
export class ProfileService {
  private apiUrl = 'http://localhost:8080/api/stakeholders/profiles/me';
  private socialUrl = 'http://localhost:8080/api/social/users';

  constructor(private http: HttpClient, private authService: AuthService) { }

  private getAuthHeaders(): HttpHeaders {
    const token = this.authService.getToken();
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }

  getProfile(): Observable<Profile> {
    const headers = this.getAuthHeaders();
    return this.http.get<Profile>(this.apiUrl, { headers });
  }

  updateProfile(data: UpdateProfileRequest): Observable<Profile> {
    const headers = this.getAuthHeaders();
    return this.http.put<Profile>(this.apiUrl, data, { headers });
  }

  getStats(userId: string): Observable<ProfileStats> {
    const headers = this.getAuthHeaders();
    return this.http.get<ProfileStats>(`${this.socialUrl}/${userId}/stats`, { headers });
  }
}