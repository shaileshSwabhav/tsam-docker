import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { FacultyBatchComponent } from './faculty-batch.component';

describe('FacultyBatchComponent', () => {
  let component: FacultyBatchComponent;
  let fixture: ComponentFixture<FacultyBatchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ FacultyBatchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(FacultyBatchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
