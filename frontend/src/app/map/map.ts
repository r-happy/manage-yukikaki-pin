import {
  Component,
  OnInit,
  AfterViewInit,
  OnDestroy,
  Inject,
  PLATFORM_ID,
} from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';

import type { Map as LeafletMap, Marker, Icon } from 'leaflet';
import { environment } from '../../environments/environment';
import { Pin } from '../types/pin.type';

@Component({
  selector: 'app-map',
  templateUrl: './map.html',
  styleUrls: ['./map.scss'], // styleUrl -> styleUrls (配列)
  imports: [RouterModule],
})
export class Map implements OnInit, AfterViewInit, OnDestroy {
  private L: any;
  private map: LeafletMap | undefined;
  private markers: Marker[] = [];
  pins: Pin[] = []; // 初期化

  public groupId: string | null = null;
  private isManualCenter = false; // 手動で中心を設定したかどうかのフラグ

  constructor(
    @Inject(PLATFORM_ID) private platformId: Object,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      const newGroupId = params['groupId'] || null;

      // グループが変更された場合、手動中心設定フラグをリセット
      if (this.groupId !== newGroupId) {
        this.isManualCenter = false;
      }

      this.groupId = newGroupId;

      // 座標が指定されている場合は保存
      const lat = params['lat'];
      const lng = params['lng'];
      const pinId = params['pinId'];

      this.fetchPins(localStorage.getItem('token') || '');

      if (this.map && lat && lng) {
        // 地図の中心を更新
        this.centerMapOnPin(parseFloat(lat), parseFloat(lng), pinId);
      }
    });
  }

  async ngAfterViewInit(): Promise<void> {
    if (isPlatformBrowser(this.platformId)) {
      this.L = await import('leaflet');
      this.initMap();
    }
  }

  ngOnDestroy(): void {
    if (this.map) {
      this.map.remove();
    }
  }

  private fetchPins(token: string) {
    if (!this.groupId) {
      this.pins = [];
      this.updateMap();
      return;
    }
    fetch(`${environment.backendUrl}/api/groups/${this.groupId}/pins`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
      })
      .then((data) => {
        this.pins = Array.isArray(data) ? data : [];
        this.updateMap();
      })
      .catch((error) => {
        this.router.navigate(['/dashboard']);
        this.pins = [];
        this.updateMap();
      });
  }

  private initMap(): void {
    if (!this.L) return;

    const iconDefault = this.L.icon({
      iconRetinaUrl: 'assets/marker-icon-2x.png',
      iconUrl: 'assets/marker-icon.png',
      shadowUrl: 'assets/marker-shadow.png',
      iconSize: [25, 41],
      iconAnchor: [12, 41],
      popupAnchor: [1, -34],
      tooltipAnchor: [16, -28],
      shadowSize: [41, 41],
    });
    this.L.Marker.prototype.options.icon = iconDefault;

    this.map = this.L.map('map');

    this.L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
    }).addTo(this.map);

    this.updateMap();
  }

  private updateMap(): void {
    if (!this.map || !this.L) {
      return;
    }

    // Clear existing markers
    this.markers.forEach((marker) => marker.remove());
    this.markers = [];

    if (this.pins && this.pins.length > 0) {
      this.pins.forEach((pin) => {
        const marker = this.L.marker([pin.latitude, pin.longitude])
          .addTo(this.map)
          .bindPopup(`<b>${pin.pin_name}</b><br>${pin.pin_type.pin_type_name}`);
        this.markers.push(marker);
      });

      // 手動で中心を設定していない場合のみfitBoundsを実行
      if (!this.isManualCenter) {
        const group = this.L.featureGroup(this.markers);
        this.map.fitBounds(group.getBounds().pad(0.1));
      }
    } else if (!this.isManualCenter) {
      // If no pins and not manually centered, set a default view
      this.map.setView([35.681236, 139.767125], 10); // Tokyo station
    }
  }

  private centerMapOnPin(lat: number, lng: number, pinId?: string): void {
    if (!this.map || !this.L) {
      return;
    }

    // 手動で中心を設定したことを記録
    this.isManualCenter = true;

    // 指定された座標に地図の中心を移動（ズームレベル18で詳細表示）
    this.map.setView([lat, lng], 18);

    // 該当するピンのポップアップを開く
    if (pinId) {
      const targetMarker = this.markers.find((marker) => {
        const pin = this.pins.find((p) => p.pin_id === pinId);
        if (pin) {
          const markerLatLng = marker.getLatLng();
          return (
            markerLatLng.lat === pin.latitude &&
            markerLatLng.lng === pin.longitude
          );
        }
        return false;
      });

      if (targetMarker) {
        targetMarker.openPopup();
      }
    }
  }
}
