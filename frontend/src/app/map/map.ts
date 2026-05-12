import {
  Component,
  OnInit,
  AfterViewInit,
  OnDestroy,
  Inject,
  PLATFORM_ID,
} from '@angular/core';
import { CommonModule, isPlatformBrowser } from '@angular/common';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';

import type {
  Map as LeafletMap,
  Marker,
  Icon,
  LeafletMouseEvent,
} from 'leaflet';
import { Subscription } from 'rxjs';
import { environment } from '../../environments/environment';
import { Pin } from '../types/pin.type';
import {
  CoordinateSelectionService,
  SelectedCoordinates,
} from '../services/coordinate-selection.service';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-map',
  templateUrl: './map.html',
  styleUrls: ['./map.scss'], // styleUrl -> styleUrls (配列)
  imports: [RouterModule, CommonModule, MatIconModule],
})
export class Map implements OnInit, AfterViewInit, OnDestroy {
  private L: any;
  private map: LeafletMap | undefined;
  private markers: Marker[] = [];
  private selectionMarker: Marker | null = null;
  private selectionSubscription?: Subscription;
  private pendingSelection: { lat: number; lng: number } | null = null;
  selectedCoordinates: SelectedCoordinates | null = null;
  pins: Pin[] = []; // 初期化

  public groupId: string | null = null;
  private isManualCenter = false; // 手動で中心を設定したかどうかのフラグ

  constructor(
    @Inject(PLATFORM_ID) private platformId: Object,
    private route: ActivatedRoute,
    private router: Router,
    private coordinateSelection: CoordinateSelectionService
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

    this.selectionSubscription =
      this.coordinateSelection.coordinates$.subscribe((coords) =>
        this.handleCoordinateSelection(coords)
      );
  }

  async ngAfterViewInit(): Promise<void> {
    if (isPlatformBrowser(this.platformId)) {
      const leafletModule = await import('leaflet');
      this.L = leafletModule.default ?? leafletModule;
      this.initMap();
    }
  }

  ngOnDestroy(): void {
    if (this.map) {
      this.map.remove();
    }
    this.selectionSubscription?.unsubscribe();
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

    const mapInstance = this.map;
    if (!mapInstance) {
      return;
    }

    mapInstance.on('click', (event: LeafletMouseEvent) => {
      const { lat, lng } = event.latlng;
      this.isManualCenter = true;
      this.setSelectionMarker(lat, lng, true);
      this.selectedCoordinates = { latitude: lat, longitude: lng };
      this.coordinateSelection.setCoordinates({
        latitude: lat,
        longitude: lng,
      });
    });

    this.updateMap();

    // The map is rendered inside a routed layout, so size can be wrong on first paint.
    queueMicrotask(() => this.map?.invalidateSize());
    setTimeout(() => this.map?.invalidateSize(), 0);
    setTimeout(() => this.map?.invalidateSize(), 200);

    if (this.pendingSelection) {
      this.setSelectionMarker(
        this.pendingSelection.lat,
        this.pendingSelection.lng,
        false
      );
      this.pendingSelection = null;
    }
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

  private handleCoordinateSelection(coords: SelectedCoordinates | null) {
    this.selectedCoordinates = coords;

    if (!coords) {
      this.pendingSelection = null;
      this.clearSelectionMarker();
      return;
    }

    const { latitude, longitude } = coords;

    if (!this.map || !this.L) {
      this.pendingSelection = { lat: latitude, lng: longitude };
      return;
    }

    this.setSelectionMarker(latitude, longitude, false);
  }

  private setSelectionMarker(lat: number, lng: number, focus = false) {
    if (!this.map || !this.L) {
      this.pendingSelection = { lat, lng };
      return;
    }

    if (!this.selectionMarker) {
      const marker = this.L.marker([lat, lng]).addTo(this.map);
      marker.bindPopup('新しいピン候補');
      this.selectionMarker = marker;
    } else {
      this.selectionMarker.setLatLng([lat, lng]);
    }

    if (focus) {
      const map = this.map;
      const marker = this.selectionMarker;

      if (map && marker) {
        const zoom = map.getZoom() < 16 ? 16 : map.getZoom();
        map.setView([lat, lng], zoom);
        marker.openPopup();
      }
    }
  }

  private clearSelectionMarker() {
    if (this.selectionMarker) {
      this.selectionMarker.remove();
      this.selectionMarker = null;
    }
  }

  clearSelectionFromMap() {
    this.coordinateSelection.clear();
  }

  get selectedCoordinateText(): string | null {
    if (!this.selectedCoordinates) {
      return null;
    }
    const { latitude, longitude } = this.selectedCoordinates;
    return `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`;
  }
}
