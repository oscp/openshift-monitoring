import {Component} from '@angular/core';

@Component({
    selector: 'app-root',
    template: `<br/><div class="container-fluid">
    <h3>OpenShift - Updatemonitoring</h3>
    <app-deamon-overview></app-deamon-overview>
    <simple-notifications [options]="notificationOptions"></simple-notifications>
</div>
    `
})
export class AppComponent {
    private notificationOptions = {
        position: ["top", "right"],
        timeOut: 3000,
        showProgressBar: true
    }
}
