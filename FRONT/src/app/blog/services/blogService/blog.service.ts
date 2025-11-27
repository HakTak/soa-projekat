import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Blog } from '../../model/blog';

@Injectable({
  providedIn: 'root'
})
export class BlogService {

  private baseUrl = 'http://localhost:5251/blog/post';

  constructor(private http: HttpClient){}

  getAllPosts(): Observable<Blog[]> {
    return this.http.get<Blog[]>(`${this.baseUrl}/all`);
  }

  toggleLike(postId: string, userId: string): Observable<Blog> {
    return this.http.post<Blog>(`${this.baseUrl}/toggleLike`, {
      postId,
      userId
    }, {headers: this.getAuthHeaders()});
  }

  createPost(formData: FormData): Observable<Blog> {
    return this.http.post<Blog>(`${this.baseUrl}/create`, formData ,  {headers: this.getAuthHeaders()});
  }

   getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt');;
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }
}
