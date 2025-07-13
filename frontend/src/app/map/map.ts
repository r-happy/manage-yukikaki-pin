import {
  Component,
  OnInit,
  AfterViewInit,
  OnDestroy,
  Inject,
  PLATFORM_ID,
} from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { ActivatedRoute } from '@angular/router';

import type { Map as LeafletMap, Marker, Icon } from 'leaflet';
import { environment } from '../../environments/environment';
import { Pin } from '../types/pin.type';

@Component({
  selector: 'app-map',
  templateUrl: './map.html',
  styleUrls: ['./map.scss'], // styleUrl -> styleUrls (配列)
})
export class Map implements OnInit, AfterViewInit, OnDestroy {
  private L: any;
  private map: LeafletMap | undefined;
  private markers: Marker[] = [];
  pins: Pin[] = []; // 初期化

  public groupId: string | null = null;

  constructor(
    @Inject(PLATFORM_ID) private platformId: Object,
    private route: ActivatedRoute
  ) {}

  ngOnInit(): void {
    this.route.queryParams.subscribe((params) => {
      this.groupId = params['groupId'] || null;
      this.fetchPins(localStorage.getItem('token') || '');
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
        console.error('Error fetching pins:', error);
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

      const group = this.L.featureGroup(this.markers);
      this.map.fitBounds(group.getBounds().pad(0.1));
    } else {
      // If no pins, set a default view
      this.map.setView([35.681236, 139.767125], 10); // Tokyo station
    }
  }
}
