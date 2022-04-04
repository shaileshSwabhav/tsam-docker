import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TalentMasterComponent } from './talent-master.component';

describe('TalentMasterComponent', () => {
  let component: TalentMasterComponent;
  let fixture: ComponentFixture<TalentMasterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TalentMasterComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TalentMasterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
