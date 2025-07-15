import { Component, OnInit, Inject, PLATFORM_ID } from '@angular/core';
import {
  Router,
  RouterModule,
  RouterOutlet,
  NavigationEnd,
} from '@angular/router';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { environment } from '../../environments/environment';
import { GroupDetail } from '../types/group-detail.type';
import { Pin, PinType } from '../types/pin.type';
import { filter } from 'rxjs/operators';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-dashboard',
  imports: [
    RouterModule,
    RouterOutlet,
    CommonModule,
    FormsModule,
    MatIconModule,
  ],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.scss',
})
export class Dashboard implements OnInit {
  groups: GroupDetail[] = [];
  pins: Pin[] = [];
  pinTypes: PinType[] = [];
  selectedGroupId: string = '';
  loading = true;
  pinsLoading = false;
  pinTypesLoading = false;
  error: string | null = null;
  pinsError: string | null = null;
  pinTypesError: string | null = null;
  showSelectGroup = true;
  showPanel = true;

  togglePanel() {
    this.showPanel = !this.showPanel;
  }

  constructor(
    private router: Router,
    @Inject(PLATFORM_ID) private platformId: Object
  ) {}

  ngOnInit() {
    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe((event: NavigationEnd) => {
        this.showSelectGroup = !event.url.includes('/profile');

        const urlParams = new URLSearchParams(event.url.split('?')[1] || '');
        const groupIdFromUrl = urlParams.get('groupId');

        if (groupIdFromUrl && groupIdFromUrl !== this.selectedGroupId) {
          this.selectedGroupId = groupIdFromUrl;
          if (isPlatformBrowser(this.platformId)) {
            const token = localStorage.getItem('token');
            if (token) {
              this.fetchPins(groupIdFromUrl);
              this.fetchPinTypes(groupIdFromUrl);
            }
          }
        } else if (!groupIdFromUrl && this.selectedGroupId !== '') {
          this.selectedGroupId = '';
          this.pins = [];
          this.pinTypes = [];
          this.pinsError = null;
          this.pinTypesError = null;
        }
      });

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
          const urlParams = new URLSearchParams(window.location.search);
          const groupIdFromUrl = urlParams.get('groupId');
          if (
            groupIdFromUrl &&
            this.groups.some((g) => g.group.group_id === groupIdFromUrl)
          ) {
            this.selectedGroupId = groupIdFromUrl;
            this.fetchPins(groupIdFromUrl);
            this.fetchPinTypes(groupIdFromUrl);
          }
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
    this.selectedGroupId = value;

    if (value) {
      this.router.navigate(['/dashboard/map'], {
        queryParams: { groupId: value },
      });
      this.fetchPins(value);
      this.fetchPinTypes(value);
    } else {
      this.router.navigate(['/dashboard/map']);
      this.pins = [];
      this.pinTypes = [];
      this.pinsError = null;
      this.pinTypesError = null;
    }
  }

  onPinClick(pin: Pin) {
    this.router.navigate(['/dashboard/map'], {
      queryParams: {
        groupId: this.selectedGroupId,
        lat: pin.latitude,
        lng: pin.longitude,
        pinId: pin.pin_id,
      },
    });
  }

  private fetchPins(groupId: string) {
    if (!isPlatformBrowser(this.platformId)) return;

    const token = localStorage.getItem('token');
    if (!token) return;

    this.pinsLoading = true;
    this.pinsError = null;

    fetch(`${environment.backendUrl}/api/groups/${groupId}/pins`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (res.ok) {
          this.pins = await res.json();
        } else if (res.status === 401) {
          localStorage.removeItem('token');
          this.router.navigate(['/signin']);
        } else {
          this.pinsError = 'ピンの取得に失敗しました';
        }
      })
      .catch(() => (this.pinsError = '通信エラー'))
      .finally(() => (this.pinsLoading = false));
  }

  private fetchPinTypes(groupId: string) {
    if (!isPlatformBrowser(this.platformId)) return;
    const token = localStorage.getItem('token');
    if (!token) return;
    this.pinTypesLoading = true;
    this.pinTypesError = null;
    fetch(`${environment.backendUrl}/api/groups/${groupId}/pin-types`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(async (res) => {
        if (res.ok) {
          this.pinTypes = await res.json();
        } else if (res.status === 401) {
          localStorage.removeItem('token');
          this.router.navigate(['/signin']);
        } else {
          this.pinTypesError = 'PinTypeの取得に失敗しました';
        }
      })
      .catch(() => (this.pinTypesError = '通信エラー'))
      .finally(() => (this.pinTypesLoading = false));
  }
}
