import { Component, ViewChild, ElementRef } from '@angular/core';
import { Comment } from '../model/comment';
import { Blog } from '../model/blog';
import { CommonModule} from '@angular/common';
import { FormsModule } from '@angular/forms';
import { BlogService } from '../services/blogService/blog.service';
import { CommentService } from '../services/commentService/comment.service';
import { AuthService } from '../../infrastructure/auth.service';
import { User} from '../../models/user.model';
import { Router } from '@angular/router';

@Component({
  selector: 'app-blogs',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './blogs.component.html',
  styleUrl: './blogs.component.css'
})
export class BlogsComponent {
  
     blogs: Blog[] = [
    {
      id: "1",
      title: "Beautiful Mountains",
      text: "Exploring the Alps during winter season...",
      tags: "",
      authorName: "John Doe",
      createdAt: new Date(),
      imageUrls: [
      ],
      comments: [
        { id: "c1", postId: "1", authorName: "Alice", text: "Amazing view!", createdAt: new Date() , updatedAt: null },
      ],
      likeCount: 0
    },
    {
      id: "2",
      authorName: "Jane Smith",
      title: "City Adventures",
      text: "Walking through the old streets...",
      tags: "",
      createdAt: new Date(),
      imageUrls: [
      ],
      comments: [],
      likeCount: 0
    }
  ];

  selectedBlog: Blog | null = null;
  newComment = '';
  editingComment: string | null = null;
  editingText: string = '';
  user: User | null;
  @ViewChild('feedContainer') feedContainer!: ElementRef;

  constructor(private blogService: BlogService, private commentService: CommentService, private authService: AuthService,
    private router: Router
  ) {
    this.blogService.getAllBlogs().subscribe({
      next: (data: Blog[]) => {
        this.blogs = data;
        console.log('Blogs fetched:', data);
      },
      error: (err) => {
        console.error('Error fetching blogs:', err);
      }
    });
    this.user = this.loadCurrentUser(); 
  }

  goToAuthorPage(authorName: string)
  {
    this.router.navigate(['/profile', authorName])
  }

  loadCurrentUser() {
   const token = localStorage.getItem('jwt');
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));

      const user: User = {
        id: payload.id,
        username: payload.username,
        email: payload.email,
        role: payload.role,
        blocked: false
      };
        return user;
    } catch (err) {
      console.error("Invalid JWT:", err);
      return null;
    }
  }

  ngAfterViewInit() {
    this.feedContainer.nativeElement.addEventListener('scroll', () => {
      // Close comments when scrolling
      if (this.selectedBlog) {
        this.selectedBlog = null;
        this.newComment = '';
      }
    });
  }

  openComments(blog: Blog) {
    if (!this.selectedBlog) 
    {
      this.selectedBlog = blog;
    }else
    {
      this.newComment = '';
      this.selectedBlog = null;
    }
  }

  addComment(blog: Blog) {
    if (!this.newComment.trim() || !this.user) return;
    const comment: Comment = {
      id: '',
      postId: blog.id,
      authorName: this.user.username,
      text: this.newComment.trim(),
      createdAt: new Date(),
      updatedAt: null
    }
    this.newComment = '';
    this.commentService.createComment(comment).subscribe({
      next: (res: Comment) => {
        console.log('Comment created:', res);
        blog.comments.push(res);
      },
      error: (err) => {
        console.error('toggleLike error:', err);
      }
    })
  }

  likeBlog(blog: Blog) {
    this.blogService.toggleLike(blog.id).subscribe({
      next: (res: Blog) => {
        if (this.selectedBlog) {
          this.selectedBlog.likeCount = res.likeCount;
        }
        console.log('toggleLike response:', res);
      },
      error: (err) => {
        console.error('toggleLike error:', err);
      }
   });
  }

  startEdit(comment: Comment) {
  this.editingComment = comment.id;
  this.editingText = comment.text;
}

  cancelEdit() {
    this.editingComment = null;
    this.editingText = '';
  }

  saveComment(comment: Comment) {
    if (!this.editingText.trim()) return;

    // Call API to update
    this.commentService.updateComment(comment).subscribe({
      next: (res) => {
        comment.text = this.editingText;
        comment.updatedAt = new Date();
        console.log('Comment updated', res)
      },
      error: (err) => console.error('Update error', err)
    });

    this.cancelEdit();
  }

  deleteComment(comment: Comment, blog: Blog) {

    this.commentService.deleteComment(comment.id).subscribe({
      next: () => {
        blog.comments = blog.comments.filter(c => c.id !== comment.id);
        console.log('Comment deleted');
      },
      error: (err) => console.error('Delete error', err)
    });
  }
}