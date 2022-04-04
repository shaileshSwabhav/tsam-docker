import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BatchModulesComponent } from './batch-modules.component';

describe('BatchModulesComponent', () => {
  let component: BatchModulesComponent;
  let fixture: ComponentFixture<BatchModulesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BatchModulesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BatchModulesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
