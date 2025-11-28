import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Blog } from '../../model/blog';

@Injectable({
  providedIn: 'root'
})
export class BlogService {

  private baseUrl = 'http://localhost:8080/api/blog/posts';

  constructor(private http: HttpClient){}

  getAllBlogs(): Observable<Blog[]> {
    return this.http.get<Blog[]>(`${this.baseUrl}`, {headers: this.getAuthHeaders()});
  }

  toggleLike(blogId: string): Observable<Blog> {
    return this.http.post<Blog>(`${this.baseUrl}/toggleLike${blogId}`, {headers: this.getAuthHeaders()});
  }

  createBlog(blog: Blog): Observable<Blog> {
    return this.http.post<Blog>(`${this.baseUrl}`, blog ,  {headers: this.getAuthHeaders()});
  }

   getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt');;
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }
}
