import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

export interface SelectedCoordinates {
  latitude: number;
  longitude: number;
}

@Injectable({ providedIn: 'root' })
export class CoordinateSelectionService {
  private readonly coordinatesSubject =
    new BehaviorSubject<SelectedCoordinates | null>(null);

  readonly coordinates$ = this.coordinatesSubject.asObservable();

  setCoordinates(coordinates: SelectedCoordinates) {
    this.coordinatesSubject.next(coordinates);
  }

  clear() {
    this.coordinatesSubject.next(null);
  }
}
