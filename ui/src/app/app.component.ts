import {Component} from '@angular/core';

@Component({
    selector: 'app-root',
    template: `
<br/><div class='container-fluid'>
    <h3>OpenShift - Updatemonitoring</h3>
    <simple-notifications [options]='notificationOptions'></simple-notifications>
    <app-deamon-overview></app-deamon-overview>
    <app-checks></app-checks>
</div>
    `
})
export class AppComponent {
    notificationOptions = {
        position: ['top', 'right'],
        timeOut: 3000,
        showProgressBar: true
    };
}
