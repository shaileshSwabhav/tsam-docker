import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgrammingProjectComponent } from './programming-project.component';

describe('ProgrammingProjectComponent', () => {
  let component: ProgrammingProjectComponent;
  let fixture: ComponentFixture<ProgrammingProjectComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProgrammingProjectComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgrammingProjectComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
