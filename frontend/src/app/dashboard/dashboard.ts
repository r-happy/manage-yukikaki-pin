import {
  Component,
  OnInit,
  Inject,
  PLATFORM_ID,
  OnDestroy,
} from '@angular/core';
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
import { User } from '../types/user.type';
import { filter } from 'rxjs/operators';
import { Subscription } from 'rxjs';
import { MatIconModule } from '@angular/material/icon';
import {
  CoordinateSelectionService,
  SelectedCoordinates,
} from '../services/coordinate-selection.service';

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
export class Dashboard implements OnInit, OnDestroy {
  groups: GroupDetail[] = [];
  pins: Pin[] = [];
  pinTypes: PinType[] = [];
  currentUserId: string | null = null;
  selectedGroupId: string = '';
  loading = true;
  pinsLoading = false;
  pinTypesLoading = false;
  error: string | null = null;
  pinsError: string | null = null;
  pinTypesError: string | null = null;
  showSelectGroup = true;
  showPanel = true;
  showAddGroupForm = false;
  addGroupLoading = false;
  addGroupError: string | null = null;
  addGroupSuccess: string | null = null;
  newGroup = {
    group_name: '',
    group_description: '',
    user_ids: '',
  };
  showAddPinTypeForm = false;
  addPinTypeLoading = false;
  addPinTypeError: string | null = null;
  addPinTypeSuccess: string | null = null;
  newPinType = {
    pin_type_name: '',
    pin_type_description: '',
  };
  showAddPinForm = false;
  addPinLoading = false;
  addPinError: string | null = null;
  addPinSuccess: string | null = null;
  newPin = {
    pin_type_id: '',
    pin_name: '',
    latitude: '',
    longitude: '',
  };
  private coordinateSubscription?: Subscription;

  togglePanel() {
    this.showPanel = !this.showPanel;
  }

  toggleAddGroupForm() {
    this.showAddGroupForm = !this.showAddGroupForm;
    if (!this.showAddGroupForm) {
      this.resetAddGroupFeedback();
    }
  }

  private resetAddGroupFeedback() {
    this.addGroupError = null;
    this.addGroupSuccess = null;
  }

  toggleAddPinTypeForm() {
    this.showAddPinTypeForm = !this.showAddPinTypeForm;
    if (!this.showAddPinTypeForm) {
      this.resetAddPinTypeFeedback();
    }
  }

  private resetAddPinTypeFeedback() {
    this.addPinTypeError = null;
    this.addPinTypeSuccess = null;
  }

  toggleAddPinForm() {
    this.showAddPinForm = !this.showAddPinForm;
    if (!this.showAddPinForm) {
      this.resetAddPinFeedback();
    }
  }

  private resetAddPinFeedback() {
    this.addPinError = null;
    this.addPinSuccess = null;
  }

  constructor(
    private router: Router,
    @Inject(PLATFORM_ID) private platformId: Object,
    private coordinateSelection: CoordinateSelectionService
  ) {}

  ngOnInit() {
    this.coordinateSubscription =
      this.coordinateSelection.coordinates$.subscribe((coords) =>
        this.applySelectedCoordinates(coords)
      );

    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe((event: NavigationEnd) => {
        this.showSelectGroup = !event.url.includes('/profile');

        const urlParams = new URLSearchParams(event.url.split('?')[1] || '');
        const groupIdFromUrl = urlParams.get('groupId');

        if (groupIdFromUrl && groupIdFromUrl !== this.selectedGroupId) {
          this.selectedGroupId = groupIdFromUrl;
          this.coordinateSelection.clear();
          if (isPlatformBrowser(this.platformId)) {
            const token = localStorage.getItem('token');
            if (token) {
              this.fetchPins(groupIdFromUrl);
              this.fetchPinTypes(groupIdFromUrl);
            }
          }
        } else if (!groupIdFromUrl && this.selectedGroupId !== '') {
          this.selectedGroupId = '';
          this.coordinateSelection.clear();
          this.pins = [];
          this.pinTypes = [];
          this.pinsError = null;
          this.pinTypesError = null;
        }
      });

    if (isPlatformBrowser(this.platformId)) {
      const token = localStorage.getItem('token');
      if (token) {
        this.ensureCurrentUserId(token);
        this.fetchGroups(token);
      } else {
        this.router.navigate(['/signin']);
      }
    }
  }

  ngOnDestroy(): void {
    this.coordinateSubscription?.unsubscribe();
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

  async createGroup() {
    if (!this.newGroup.group_name || !this.newGroup.group_description) {
      this.addGroupError = 'グループ名と説明は必須です。';
      this.addGroupSuccess = null;
      return;
    }

    if (!isPlatformBrowser(this.platformId)) return;
    const token = localStorage.getItem('token');
    if (!token) {
      this.router.navigate(['/signin']);
      return;
    }

    this.addGroupLoading = true;
    this.resetAddGroupFeedback();

    const currentUserId = await this.ensureCurrentUserId(token);
    if (!currentUserId) {
      this.addGroupError =
        'ユーザー情報の取得に失敗しました。再度サインインしてください。';
      this.addGroupLoading = false;
      return;
    }

    const params = new URLSearchParams();
    params.set('group_name', this.newGroup.group_name.trim());
    params.set('group_description', this.newGroup.group_description.trim());
    const additionalUserIds = this.newGroup.user_ids
      .split(',')
      .map((id) => id.trim())
      .filter((id) => id.length > 0 && id !== currentUserId);

    const userIdSet = new Set<string>([currentUserId, ...additionalUserIds]);
    params.set('user_ids', Array.from(userIdSet).join(','));

    try {
      const res = await fetch(`${environment.backendUrl}/api/groups`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
        },
        body: params.toString(),
      });

      if (!res.ok) {
        if (res.status === 401) {
          localStorage.removeItem('token');
          this.router.navigate(['/signin']);
          return;
        }
        const message = await res.text();
        this.addGroupError = message || 'グループの作成に失敗しました。';
        return;
      }

      this.addGroupSuccess = 'グループを作成しました。';
      this.newGroup = {
        group_name: '',
        group_description: '',
        user_ids: '',
      };
      this.fetchGroups(token);
    } catch (error) {
      this.addGroupError = '通信エラーが発生しました。';
    } finally {
      this.addGroupLoading = false;
    }
  }

  onGroupChange(event: Event) {
    const value = (event.target as HTMLSelectElement).value;
    this.selectedGroupId = value;
    this.coordinateSelection.clear();

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

  async createPinType() {
    if (!this.selectedGroupId) {
      this.addPinTypeError = '先にグループを選択してください。';
      this.addPinTypeSuccess = null;
      return;
    }

    if (
      !this.newPinType.pin_type_name ||
      !this.newPinType.pin_type_description
    ) {
      this.addPinTypeError = '名称と説明は必須です。';
      this.addPinTypeSuccess = null;
      return;
    }

    if (!isPlatformBrowser(this.platformId)) return;
    const token = localStorage.getItem('token');
    if (!token) {
      this.router.navigate(['/signin']);
      return;
    }

    this.addPinTypeLoading = true;
    this.resetAddPinTypeFeedback();

    const params = new URLSearchParams();
    params.set('pin_type_name', this.newPinType.pin_type_name.trim());
    params.set(
      'pin_type_description',
      this.newPinType.pin_type_description.trim()
    );

    try {
      const res = await fetch(
        `${environment.backendUrl}/api/groups/${this.selectedGroupId}/pin-types`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
          },
          body: params.toString(),
        }
      );

      if (!res.ok) {
        if (res.status === 401) {
          localStorage.removeItem('token');
          this.router.navigate(['/signin']);
          return;
        }
        const message = await res.text();
        this.addPinTypeError = message || 'PinTypeの作成に失敗しました。';
        return;
      }

      this.addPinTypeSuccess = 'PinTypeを追加しました。';
      this.newPinType = {
        pin_type_name: '',
        pin_type_description: '',
      };
      this.fetchPinTypes(this.selectedGroupId);
    } catch (error) {
      this.addPinTypeError = '通信エラーが発生しました。';
    } finally {
      this.addPinTypeLoading = false;
    }
  }

  async createPin() {
    if (!this.selectedGroupId) {
      this.addPinError = 'グループを選択してください。';
      this.addPinSuccess = null;
      return;
    }

    if (
      !this.newPin.pin_type_id ||
      !this.newPin.pin_name ||
      !this.newPin.latitude ||
      !this.newPin.longitude
    ) {
      this.addPinError = 'すべての項目を入力してください。';
      this.addPinSuccess = null;
      return;
    }

    const latitude = parseFloat(this.newPin.latitude);
    const longitude = parseFloat(this.newPin.longitude);

    if (Number.isNaN(latitude) || Number.isNaN(longitude)) {
      this.addPinError = '緯度・経度は数値で入力してください。';
      this.addPinSuccess = null;
      return;
    }

    if (!isPlatformBrowser(this.platformId)) return;
    const token = localStorage.getItem('token');
    if (!token) {
      this.router.navigate(['/signin']);
      return;
    }

    this.addPinLoading = true;
    this.resetAddPinFeedback();

    const params = new URLSearchParams();
    params.set('pin_type_id', this.newPin.pin_type_id);
    params.set('pin_name', this.newPin.pin_name.trim());
    params.set('latitude', latitude.toString());
    params.set('longitude', longitude.toString());

    try {
      const res = await fetch(
        `${environment.backendUrl}/api/groups/${this.selectedGroupId}/pins`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8',
          },
          body: params.toString(),
        }
      );

      if (!res.ok) {
        if (res.status === 401) {
          localStorage.removeItem('token');
          this.router.navigate(['/signin']);
          return;
        }
        const message = await res.text();
        this.addPinError = message || 'ピンの追加に失敗しました。';
        return;
      }

      this.addPinSuccess = 'ピンを追加しました。';
      this.coordinateSelection.clear();
      this.newPin = {
        pin_type_id: '',
        pin_name: '',
        latitude: '',
        longitude: '',
      };
      this.fetchPins(this.selectedGroupId);
    } catch (error) {
      this.addPinError = '通信エラーが発生しました。';
    } finally {
      this.addPinLoading = false;
    }
  }

  private async ensureCurrentUserId(token: string): Promise<string | null> {
    if (this.currentUserId) {
      return this.currentUserId;
    }

    try {
      const res = await fetch(`${environment.backendUrl}/api/me`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      if (res.ok) {
        const user: User = await res.json();
        this.currentUserId = user.user_id;
        return this.currentUserId;
      }

      if (res.status === 401) {
        localStorage.removeItem('token');
        this.router.navigate(['/signin']);
        return null;
      }

      return null;
    } catch (error) {
      return null;
    }
  }

  private applySelectedCoordinates(coords: SelectedCoordinates | null) {
    if (!coords) {
      this.newPin.latitude = '';
      this.newPin.longitude = '';
      return;
    }

    const latitude = coords.latitude.toFixed(6);
    const longitude = coords.longitude.toFixed(6);

    this.newPin.latitude = latitude;
    this.newPin.longitude = longitude;
    this.showAddPinForm = true;
    this.addPinError = null;
  }

  clearSelectedCoordinates() {
    this.coordinateSelection.clear();
    this.newPin.latitude = '';
    this.newPin.longitude = '';
  }
}
