import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentDashboardComponent } from './talent-dashboard.component';

describe('TalentDashboardComponent', () => {
  let component: TalentDashboardComponent;
  let fixture: ComponentFixture<TalentDashboardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentDashboardComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentDashboardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
