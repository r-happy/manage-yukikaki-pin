import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { environment } from '../../environments/environment';
import { Pin as PinType, PinType as PinTypeData } from '../types/pin.type';

@Component({
  selector: 'app-pin',
  imports: [RouterModule, CommonModule, FormsModule],
  templateUrl: './pin.html',
  styleUrl: './pin.scss',
})
export class Pin implements OnInit {
  pinId: string | null = null;
  groupId: string | null = null;
  pin: PinType | null = null;
  pinTypes: PinTypeData[] = [];

  // 編集用フォームデータ
  editMode = false;
  pinName = '';
  selectedPinTypeId = '';
  latitude = 0;
  longitude = 0;

  errorMessage = '';
  successMessage = '';

  constructor(private router: Router, private route: ActivatedRoute) {}

  ngOnInit() {
    this.route.queryParams.subscribe((params) => {
      this.pinId = params['pinId'] || null;
      this.groupId = params['groupId'] || null;

      if (this.pinId && this.groupId) {
        this.loadPin();
        this.loadPinTypes();
      }
    });
  }

  private loadPin() {
    const token = localStorage.getItem('token');
    if (!token || !this.groupId) return;

    fetch(`${environment.backendUrl}/api/groups/${this.groupId}/pins`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => {
        if (!response.ok) throw new Error('Failed to load pin');
        return response.json();
      })
      .then((pins: PinType[]) => {
        this.pin = pins.find((p) => p.pin_id === this.pinId) || null;
        if (this.pin) {
          this.pinName = this.pin.pin_name;
          this.selectedPinTypeId = this.pin.pin_type_id;
          this.latitude = this.pin.latitude;
          this.longitude = this.pin.longitude;
        }
      })
      .catch((error) => {
        this.errorMessage = 'ピンの読み込みに失敗しました';
        console.error(error);
      });
  }

  private loadPinTypes() {
    const token = localStorage.getItem('token');
    if (!token || !this.groupId) return;

    fetch(`${environment.backendUrl}/api/groups/${this.groupId}/pin-types`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => {
        if (!response.ok) throw new Error('Failed to load pin types');
        return response.json();
      })
      .then((data: PinTypeData[]) => {
        this.pinTypes = data;
      })
      .catch((error) => {
        console.error('Failed to load pin types:', error);
      });
  }

  enableEditMode() {
    this.editMode = true;
    this.errorMessage = '';
    this.successMessage = '';
  }

  cancelEdit() {
    this.editMode = false;
    if (this.pin) {
      this.pinName = this.pin.pin_name;
      this.selectedPinTypeId = this.pin.pin_type_id;
      this.latitude = this.pin.latitude;
      this.longitude = this.pin.longitude;
    }
    this.errorMessage = '';
    this.successMessage = '';
  }

  updatePin() {
    if (!this.pinId || !this.groupId) return;

    const token = localStorage.getItem('token');
    if (!token) {
      this.errorMessage = '認証が必要です';
      return;
    }

    const formData = new FormData();
    formData.append('pin_name', this.pinName);
    formData.append('pin_type_id', this.selectedPinTypeId);
    formData.append('latitude', this.latitude.toString());
    formData.append('longitude', this.longitude.toString());

    fetch(
      `${environment.backendUrl}/api/groups/${this.groupId}/pins/${this.pinId}`,
      {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      }
    )
      .then((response) => {
        if (!response.ok) throw new Error('Failed to update pin');
        return response.json();
      })
      .then((updatedPin: PinType) => {
        this.pin = updatedPin;
        this.editMode = false;
        this.successMessage = 'ピンを更新しました';
        setTimeout(() => {
          this.successMessage = '';
          this.router.navigate(['/map'], {
            queryParams: { groupId: this.groupId },
          });
        }, 1500);
      })
      .catch((error) => {
        this.errorMessage = 'ピンの更新に失敗しました';
        console.error(error);
      });
  }

  deletePin() {
    if (!this.pinId || !this.groupId) return;

    if (!confirm('本当にこのピンを削除しますか？')) {
      return;
    }

    const token = localStorage.getItem('token');
    if (!token) {
      this.errorMessage = '認証が必要です';
      return;
    }

    fetch(
      `${environment.backendUrl}/api/groups/${this.groupId}/pins/${this.pinId}`,
      {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    )
      .then((response) => {
        if (!response.ok) throw new Error('Failed to delete pin');
        return response.json();
      })
      .then(() => {
        this.successMessage = 'ピンを削除しました';
        setTimeout(() => {
          this.router.navigate(['/map'], {
            queryParams: { groupId: this.groupId },
          });
        }, 1500);
      })
      .catch((error) => {
        this.errorMessage = 'ピンの削除に失敗しました';
        console.error(error);
      });
  }

  goBack() {
    this.router.navigate(['/map'], {
      queryParams: { groupId: this.groupId },
    });
  }
}
