import { Component } from '@angular/core';
// Add the RxJS Observable operators.
import './rxjs-operators';

@Component({
  selector: 'my-app',
  template: `
  <h1>{{title}}</h1>
  <nav>
    <a routerLink="/dashboard" routerLinkActive="active">Dashboard</a>
    <a routerLink="/heroes" routerLinkActive="active">Heroes</a>
  </nav>
  <router-outlet></router-outlet>
  `,
  styleUrls: ['app/app.component.css'],
})

export class AppComponent {
  title = 'Tour of Heroes';
}
