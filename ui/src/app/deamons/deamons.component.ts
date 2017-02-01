import {Component, OnInit} from '@angular/core';
import {SocketService} from '../socket.service';
import {SocketType} from '../shared/socket.types';
import {NotificationsService} from 'angular2-notifications';

@Component({
  selector: 'app-deamon-overview',
  templateUrl: './deamons.component.html'
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
