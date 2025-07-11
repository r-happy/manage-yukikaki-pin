import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { environment } from '../../environments/environment';
import { User } from '../types/user.type';
import { Router, RouterModule } from '@angular/router';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.html',
  styleUrl: './profile.scss',
  standalone: true,
  imports: [CommonModule, RouterModule],
})
export class Profile implements OnInit {
  user: User | null = null;
  loading = true;
  error: string | null = null;

  constructor(private router: Router) {
    // 初期化処理
    this.user = null;
    this.loading = true;
    this.error = null;
  }

  ngOnInit() {
    fetch(`${environment.backendUrl}/api/me`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token') || ''}`,
      },
    })
      .then(async (res) => {
        if (res.ok) {
          this.user = await res.json();
        } else {
          this.error = 'ユーザー情報の取得に失敗しました';
        }
      })
      .catch(() => {
        this.error = '通信エラー';
      })
      .finally(() => {
        this.loading = false;
      });
  }

  onLogout() {
    localStorage.removeItem('token');
    this.router.navigate(['/signin']);
  }
}
