import {Component, OnInit} from '@angular/core';
import {SocketService} from '../socket.service';
import {SocketType} from '../shared/socket.types';
import {NotificationsService} from 'angular2-notifications';

@Component({
  selector: 'app-daemon-overview',
  templateUrl: './daemons.component.html'
})
export class DaemonsComponent implements OnInit {
  private daemons: any;

  constructor(private socketService: SocketService, private notificationService: NotificationsService) {
    this.getDaemons();
  }

  ngOnInit() {
    this.socketService.websocket.subscribe(
      msg => {
        let data = JSON.parse(msg.data);
        switch (data.Type) {
          case SocketType.ALL_DAEMONS:
            this.daemons = data.Message.sort();
            break;
          case SocketType.NEW_DAEMON:
            this.notificationService.info('Daemon joined', 'New daemon joined: ' + data.Message);
            this.getDaemons();
            break;
          case SocketType.DAEMON_LEFT:
            this.notificationService.info('Daemon left', 'Daemon left: ' + data.Message);
            this.getDaemons();
        }
      }
    );
  }

  getDaemons() {
    this.socketService.websocket.next({Type: SocketType.ALL_DAEMONS});
  }
}
