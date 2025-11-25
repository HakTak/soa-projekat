import { Component, ViewChild, ElementRef } from '@angular/core';
import { Comment } from '../model/comment';
import { Blog } from '../model/blog';
import { CommonModule} from '@angular/common';
import { FormsModule } from '@angular/forms';
import { BlogService } from '../services/blogService/blog.service';
import { CommentService } from '../services/commentService/comment.service';

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
      content: "Exploring the Alps during winter season...",
      tags: "",
      authorName: "John Doe",
      createdAt: new Date(),
      imageUrls: [
        "https://picsum.photos/400/300?1",
        "https://picsum.photos/400/300?2"
      ],
      comments: [
        { id: "c1", postId: "1", authorName: "Alice", content: "Amazing view!", createdAt: new Date() , updatedAt: null },
      ],
      likeCount: 0
    },
    {
      id: "2",
      authorName: "Jane Smith",
      title: "City Adventures",
      content: "Walking through the old streets...",
      tags: "",
      createdAt: new Date(),
      imageUrls: [
        "https://picsum.photos/400/300?3"
      ],
      comments: [],
      likeCount: 0
    }
  ];

  selectedBlog: Blog | null = null;
  newComment = '';
  editingComment: string | null = null;
  editingContent: string = '';
  @ViewChild('feedContainer') feedContainer!: ElementRef;

  constructor(private blogService: BlogService, private commentService: CommentService) {
    this.blogService.getAllPosts().subscribe({
      next: (data: Blog[]) => {
        this.blogs = data;
        console.log('Blogs fetched:', data);
      },
      error: (err) => {
        console.error('Error fetching blogs:', err);
      }
    });
      
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
    if (!this.newComment.trim()) return;
    const comment: Comment = {
      id: '',
      postId: blog.id,
      authorName: 'You',
      content: this.newComment.trim(),
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
    blog.comments.push(comment);
  }

  likeBlog(blog: Blog) {
    blog.likeCount++;
    this.blogService.toggleLike(blog.id, 'currentUserId').subscribe({
      next: (res) => {
        console.log('toggleLike response:', res);
      },
      error: (err) => {
        console.error('toggleLike error:', err);
      }
   });
  }

  startEdit(comment: Comment) {
  this.editingComment = comment.id;
  this.editingContent = comment.content;
}

  cancelEdit() {
    this.editingComment = null;
    this.editingContent = '';
  }

  saveComment(comment: Comment) {
    if (!this.editingContent.trim()) return;

    comment.content = this.editingContent;
    comment.updatedAt = new Date();

    // Call API to update
    this.commentService.updateComment(comment).subscribe({
      next: (res) => console.log('Comment updated', res),
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