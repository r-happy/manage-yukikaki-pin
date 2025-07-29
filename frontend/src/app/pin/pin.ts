import { Component } from '@angular/core';
import { ActivatedRoute, RouterModule } from '@angular/router';

@Component({
  selector: 'app-pin',
  imports: [RouterModule],
  templateUrl: './pin.html',
  styleUrl: './pin.scss',
})
export class Pin {
  private pinId: string | null = null;

  constructor(private router: RouterModule, private route: ActivatedRoute) {}

  ngOnInit() {
    this.route.queryParams.subscribe((params) => {
      const newPinId = params['pinId'] || null;
    });
  }
}
