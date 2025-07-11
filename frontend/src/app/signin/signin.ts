import { Component } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { environment } from '../../environments/environment';
import { Router, RouterModule } from '@angular/router';

@Component({
  selector: 'app-signin',
  imports: [ReactiveFormsModule, RouterModule],
  templateUrl: './signin.html',
  styleUrl: './signin.scss',
})
export class Signin {
  signInForm: FormGroup;

  constructor(private fb: FormBuilder, private router: Router) {
    this.signInForm = this.fb.nonNullable.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required]],
    });
  }

  onSubmit() {
    if (this.signInForm.valid) {
      const formData = new FormData();
      formData.append('email', this.signInForm.value.email);
      formData.append('password', this.signInForm.value.password);

      fetch(`${environment.backendUrl}/signin`, {
        method: 'POST',
        body: formData,
      })
        .then(async (response) => {
          if (response.ok) {
            const data = await response.json();
            localStorage.setItem('token', data.token);
            this.router.navigate(['/dashboard']);
          } else {
            const error = await response.json();
            alert(error.message || 'サインインに失敗しました');
          }
        })
        .catch((error) => {
          console.error('Error:', error);
          alert('通信エラーが発生しました');
        });
    } else {
      console.log('Form is invalid');
    }
  }
}
