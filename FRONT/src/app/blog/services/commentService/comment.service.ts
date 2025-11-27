import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Comment } from '../../model/comment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CommentService {

  private baseUrl = 'http://localhost:5251/blog/comments';

  constructor(private http: HttpClient) {}

  getCommentsByPostId(postId: string):  Observable<Comment[]> {
    return this.http.get<Comment[]>(`${this.baseUrl}/getAllByPostId/${postId}`);
  }

  createComment(comment: Comment): Observable<Comment> {
    return this.http.post<Comment>(`${this.baseUrl}/create`, comment);
  }

  updateComment(comment: Comment): Observable<Comment>{
    return this.http.put<Comment>(`${this.baseUrl}/edit/`, comment);
  }

  deleteComment(commentId: string): Observable<null> {
    return this.http.delete<null>(`${this.baseUrl}/delete/${commentId}`);
  }

  getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt');;
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }
}
