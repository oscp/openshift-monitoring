import {Component, OnInit} from '@angular/core';
import {SocketService} from '../socket.service';
import {SocketType} from '../shared/socket.types';
import {NotificationsService} from 'angular2-notifications';

@Component({
  selector: 'app-deamon-overview',
  template: `
        <h4>Connected Deamons</h4>
        <table class='table table-striped'>
            <thead class='thead-inverse'>
                <tr>
                <th>Hostname</th>
                <th>Type</th>
                <th>Checks started/finished</th>
                </tr>
            </thead>
            <tbody>
                <tr *ngFor='let d of deamons'>
                    <td>{{d.Hostname}}</td>
                    <td>{{d.DeamonType}}</td>
                    <td>{{d.StartedChecks}} / {{d.FinishedChecks}}</td>
                </tr>
            </tbody>
        </table>
    `
})
export class DeamonsComponent implements OnInit {
  private deamons: any;

  constructor(private socketService: SocketService, private notificationService: NotificationsService) {
    this.getDeamons();
  }

  ngOnInit() {
    this.socketService.websocket.subscribe(
      msg => {
        let data = JSON.parse(msg.data);
        switch (data.Type) {
          case SocketType.ALL_DEAMONS:
            this.deamons = data.Message;
            break;
          case SocketType.NEW_DEAMON:
            this.notificationService.info('Deamon joined', 'New deamon joined: ' + data.Message);
            this.getDeamons();
            break;
          case SocketType.DEAMON_LEFT:
            this.notificationService.info('Deamon left', 'Deamon left: ' + data.Message);
            this.getDeamons();
        }
      }
    );
  }

  getDeamons() {
    this.socketService.websocket.next({Type: SocketType.ALL_DEAMONS});
  }
}
