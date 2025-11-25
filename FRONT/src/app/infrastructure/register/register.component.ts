import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
// Import iz foldera iznad (infrastructure)
import { AuthService, RegisterRequest } from '../auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './register.component.html',
  styleUrl: './register.component.css'
})
export class RegisterComponent {

  registerData: RegisterRequest = {
    username: '',
    password: '',
    email: '',
    role: 'TOURIST'
  };

  constructor(private authService: AuthService, private router: Router) { }

  onRegister() {
    this.authService.register(this.registerData).subscribe({
      next: (res) => {
        alert('Registracija uspesna!');
        this.router.navigate(['/login']);
      },
      error: (err) => {
        console.error(err);
        alert('Greska pri registraciji');
      }
    });
  }
}