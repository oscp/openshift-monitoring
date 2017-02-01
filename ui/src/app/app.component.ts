import {Component} from '@angular/core';

@Component({
    selector: 'app-root',
    template: `<br/>
<div class='container-fluid'>
    <h3>OpenShift - Updatemonitoring</h3>
    <simple-notifications [options]='notificationOptions'></simple-notifications>
    <div class="row">
        <div class="col">
            <app-deamon-overview></app-deamon-overview>
        </div>
        <div class="col">
            <app-checks></app-checks>
        </div>
    </div>
    <app-results></app-results>
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
