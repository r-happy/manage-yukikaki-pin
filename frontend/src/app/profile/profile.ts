import { Component, OnInit, Inject, PLATFORM_ID } from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { environment } from '../../environments/environment';
import { User } from '../types/user.type';
import { Router, RouterModule } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.html',
  styleUrl: './profile.scss',
  standalone: true,
  imports: [CommonModule, RouterModule, MatIconModule],
})
export class Profile implements OnInit {
  user: User | null = null;
  loading = true;
  error: string | null = null;

  constructor(
    private router: Router,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {
    // 初期化処理
    this.user = null;
    this.loading = true;
    this.error = null;
  }

  ngOnInit() {
    // SSR環境では何もしない
    if (!isPlatformBrowser(this.platformId)) {
      return;
    }

    // クライアントサイドでトークンチェック
    const token = localStorage.getItem('token');
    if (!token) {
      this.router.navigate(['/signin']);
      return;
    }

    // トークンがある場合のみAPIリクエスト
    this.fetchUserData(token);
  }

  private fetchUserData(token: string) {
    fetch(`${environment.backendUrl}/api/me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (res.ok) {
          this.user = await res.json();
        } else if (res.status === 401) {
          // トークンが無効な場合
          localStorage.removeItem('token');
          this.router.navigate(['/signin']);
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
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem('token');
    }
    this.router.navigate(['/signin']);
  }
}
