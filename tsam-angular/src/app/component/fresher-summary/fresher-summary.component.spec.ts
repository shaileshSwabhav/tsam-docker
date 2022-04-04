import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { FresherSummaryComponent } from './fresher-summary.component';

describe('FresherSummaryComponent', () => {
  let component: FresherSummaryComponent;
  let fixture: ComponentFixture<FresherSummaryComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ FresherSummaryComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FresherSummaryComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
