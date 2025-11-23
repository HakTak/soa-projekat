import { Component } from '@angular/core';
import { Blog } from '../model/blog';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-blog',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './blog.component.html',
  styleUrl: './blog.component.css'
})
export class BlogComponent {
  public currentBlog: Blog = {
    id: '',
    title: '',
    content: '',
    createdAt: '' as unknown as Date,
    tags: '',
    imageUrls: [],
    comments: []
    }
    currentFiles: File[] = [];

    resetForm() {
    this.currentBlog = {
      id: '',
      title: '',
      tags: '',
      content: '',
      createdAt: new Date(),
      imageUrls: [],
      comments: []
    };
    this.currentFiles = [];
  }

  onFileSelected(event: any) {
    this.currentFiles = [];
    const files: FileList = event.target.files;
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) return;
    const filesArray = Array.from(input.files);
    this.currentFiles.push(...filesArray);
  }

  publishBlog() {
    const formData = new FormData();
    this.currentBlog.createdAt = new Date();
    console.log('Publishing blog:', this.currentBlog);
    formData.append('blog', JSON.stringify(this.currentBlog));
    formData.append('images', JSON.stringify(this.currentFiles));
    this.resetForm();
  }
}
