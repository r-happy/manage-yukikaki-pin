import { Component, OnInit, Inject, PLATFORM_ID } from '@angular/core';
import { Router, RouterModule, RouterOutlet } from '@angular/router';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { environment } from '../../environments/environment';
import { GroupDetail } from '../types/group-detail.type';

@Component({
  selector: 'app-dashboard',
  imports: [RouterModule, RouterOutlet, CommonModule],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.scss',
})
export class Dashboard implements OnInit {
  groups: GroupDetail[] = [];
  loading = true;
  error: string | null = null;

  constructor(
    private router: Router,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  ngOnInit() {
    if (isPlatformBrowser(this.platformId)) {
      const token = localStorage.getItem('token');
      if (token) {
        this.fetchGroups(token);
      } else {
        this.router.navigate(['/signin']);
      }
    }
  }

  fetchGroups(token: string) {
    this.loading = true;
    fetch(`${environment.backendUrl}/api/me/groups`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (res.ok) {
          this.groups = await res.json();
        } else if (res.status === 401) {
          localStorage.removeItem('token');
          this.router.navigate(['/signin']);
        } else {
          this.error = 'グループの取得に失敗しました';
        }
      })
      .catch(() => (this.error = '通信エラー'))
      .finally(() => (this.loading = false));
  }

  onGroupChange(event: Event) {
    const value = (event.target as HTMLSelectElement).value;
    if (value) {
      this.router.navigate(['/dashboard/map'], {
        queryParams: { groupId: value },
      });
    } else {
      this.router.navigate(['/dashboard/map']);
    }
  }
}
