import { Component } from '@angular/core';
import { Blog } from '../model/blog';
import { FormsModule } from '@angular/forms';
import { BlogService } from '../services/blogService/blog.service';
import { FormGroup, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { User} from '../../models/user.model';


@Component({
  selector: 'app-blog',
  standalone: true,
  imports: [FormsModule, ReactiveFormsModule],
  templateUrl: './blog.component.html',
  styleUrl: './blog.component.css'
})
export class BlogComponent {
  public blogForm: FormGroup;
  public currentFiles: File[];
  public user: User | null;
  
  constructor(private blogService: BlogService) {
    this.blogForm = new FormGroup({
      title: new FormControl('', Validators.required),
      content: new FormControl('', Validators.required),
      tags: new FormControl('')
    }); 
    this.currentFiles = [];
    this.user = this.loadCurrentUser()
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

  resetForm() {
      this.blogForm.reset();
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

  convertFilesToBase64(files: File[]): Promise<string[]> {
    const promises = files.map(file => {
      return new Promise<string>((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => {
          // Keep the full Data URL including "data:image/...;base64,"
          resolve(reader.result as string);
        };
        reader.onerror = error => reject(error);
        reader.readAsDataURL(file); // automatically detects MIME type
      });
    });

    return Promise.all(promises);
}


  parseTags() {
    let value = this.blogForm.value.tags
    value = value.replace(/ +/g, ","); 
    value = value.replace(/,+/g, ","); 
    this.blogForm.patchValue({ tags: value });
}

  async publishBlog() {
    if (!this.user) { return ;}
    const imageUrls = await this.convertFilesToBase64(this.currentFiles);
    const currentBlog: Blog = {
      id: '',
      title: this.blogForm.value.title,
      authorName:  this.user.username,
      imageUrls: imageUrls,
      comments: [],
      likeCount: 0,
      text: this.blogForm.value.text,
      tags: this.blogForm.value.tags.replace(/^,+/, "").replace(/,+$/, ""),
      createdAt: new Date()
    }
    this.blogService.createBlog(currentBlog).subscribe({
      next: (response: Blog) => {
        console.log('Blog post created successfully:', response);
      },
      error: (error) => {
        console.error('Error creating blog:', error);
      }
    });
    this.resetForm();
  }
}
