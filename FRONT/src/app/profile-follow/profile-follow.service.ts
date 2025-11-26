import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, forkJoin, map, switchMap, of } from 'rxjs';
import { AuthService } from '../infrastructure/auth.service';

// IMPORTUJEMO MODELE IZ ZAJEDNICKOG FOLDERA
import { Profile, Stats, FullProfileDisplay, RecommendationRaw } from '../models/profile-follow.model';

@Injectable({
  providedIn: 'root'
})
export class ProfileFollowService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient, private authService: AuthService) { }


  // --- STAKEHOLDERS POZIVI (Podaci o korisniku) ---
  getProfile(userId: string): Observable<Profile> {
    const headers = this.authService.getAuthHeaders();
    return this.http.get<Profile>(`${this.apiUrl}/stakeholders/profiles/${userId}`, { headers });
  }

  // --- FOLLOWER POZIVI (Statistika i Akcije) ---

  getStats(userId: string): Observable<Stats> {
    const headers = this.authService.getAuthHeaders();
    return this.http.get<Stats>(`${this.apiUrl}/social/users/${userId}/stats`, { headers });
  }

  follow(targetId: string): Observable<any> {
    const headers = this.authService.getAuthHeaders();
    return this.http.post(`${this.apiUrl}/social/follow`, { target_id: targetId }, { headers });
  }

  unfollow(targetId: string): Observable<any> {
    const headers = this.authService.getAuthHeaders();
    return this.http.post(`${this.apiUrl}/social/unfollow`, { target_id: targetId }, { headers });
  }

  getRecommendations(): Observable<{ recommendations: RecommendationRaw[] }> {
    const headers = this.authService.getAuthHeaders();
    return this.http.get<{ recommendations: RecommendationRaw[] }>(`${this.apiUrl}/social/recommendations`, { headers });
  }

  // Provera da li pratim korisnika
  // (Zahteva da na backendu postoji POST /api/social/is-following)
  isFollowing(fromId: string, targetId: string): Observable<{ isFollowing: boolean }> {
    const headers = this.authService.getAuthHeaders();
    return this.http.post<{ isFollowing: boolean }>(`${this.apiUrl}/social/is-following`, {
      follower_id: fromId,
      followee_id: targetId,
    }, { headers });
  }

  // --- AGREGACIJA PODATAKA ---
  // Ova funkcija uzima listu preporuka (samo ID-evi) i za svakog dovlaci Ime i Sliku iz Stakeholders servisa
  getEnrichedRecommendations(): Observable<FullProfileDisplay[]> {
    return this.getRecommendations().pipe(
      switchMap((res) => {
        const recs = res.recommendations || [];

        if (recs.length === 0) return of([]);

        // Za svaki ID, pravimo paralelni zahtev ka Stakeholders i Follower servisu
        const requests = recs.map(rec =>
          forkJoin({
            profile: this.getProfile(rec.userId),
            stats: this.getStats(rec.userId)
          }).pipe(
            map(data => ({
              ...data.profile,
              stats: data.stats,
              mutuals: rec.mutualConnections
            } as FullProfileDisplay))
          )
        );

        return forkJoin(requests);
      })
    );
  }
}