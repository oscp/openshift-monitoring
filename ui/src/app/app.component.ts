import {Component} from '@angular/core';

@Component({
  selector: 'app-root',
  template: `
    <nav class="navbar is-dark">
      <div class="navbar-brand">
        <span class="icon is-large"><i class="fas fa-cloud"></i> </span>
        <a class="navbar-item" href="#">OpenShift - Updatemonitoring</a>
      </div>
    </nav>
    <simple-notifications [options]='notificationOptions'></simple-notifications>

    <section class="section">
      <div class="container">
        <div class="columns">
          <div class="column is-half">
            <app-daemon-overview></app-daemon-overview>
          </div>
          <div class="column is-half">
            <app-checks></app-checks>
          </div>
        </div>
        <app-results></app-results>
      </div>
    </section>
  `
})
export class AppComponent {
  notificationOptions = {
    position: ['top', 'right'],
    timeOut: 3000,
    showProgressBar: true,
    maxStack: 8,
    preventDuplicates: true,
    maxLength: 10
  };
}
