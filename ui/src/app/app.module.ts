import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {FormsModule} from '@angular/forms';
import {HttpModule} from '@angular/http';
import {AppComponent} from './app.component';
import {SocketService} from './socket.service';
import {SimpleNotificationsModule} from 'angular2-notifications';
import {DeamonsComponent} from './deamons/deamons.component';
import { ChecksComponent } from './checks/checks.component';
import { ResultsComponent } from './results/results.component';
import {ChartsModule} from "ng2-charts";

@NgModule({
  declarations: [
    AppComponent,
    DeamonsComponent,
    ChecksComponent,
    ResultsComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    SimpleNotificationsModule,
    ChartsModule
  ],
  providers: [SocketService],
  bootstrap: [AppComponent]
})
export class AppModule {
}
