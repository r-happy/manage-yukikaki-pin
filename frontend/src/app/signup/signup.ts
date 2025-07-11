import { Component } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { environment } from '../../environments/environment';
import { RouterModule, Router } from '@angular/router';

@Component({
  selector: 'app-signup',
  imports: [ReactiveFormsModule, RouterModule],
  templateUrl: './signup.html',
  styleUrl: './signup.scss',
})
export class Signup {
  signUpForm: FormGroup;

  constructor(private fb: FormBuilder, private router: Router) {
    this.signUpForm = this.fb.nonNullable.group({
      user_name: ['', [Validators.required]],
      user_email: ['', [Validators.required, Validators.email]],
      user_password: ['', [Validators.required]],
      user_confirmPassword: ['', [Validators.required]],
    });
  }

  onSubmit() {
    if (this.signUpForm.valid) {
      const formData = new FormData();
      formData.append('user_name', this.signUpForm.value.user_name);
      formData.append('user_email', this.signUpForm.value.user_email);
      formData.append('user_password', this.signUpForm.value.user_password);

      fetch(`${environment.backendUrl}/signup`, {
        method: 'POST',
        body: formData,
      }).then(async (response) => {
        if (response.ok) {
          this.router.navigate(['/signin']);
        } else {
          const error = await response.json();
          alert(error.message || 'サインアップに失敗しました');
        }
      });
    } else {
      console.log('Form is invalid');
    }
  }
}
