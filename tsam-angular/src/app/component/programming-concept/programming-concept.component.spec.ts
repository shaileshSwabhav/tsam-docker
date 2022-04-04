import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgrammingConceptComponent } from './programming-concept.component';

describe('ProgrammingConceptComponent', () => {
  let component: ProgrammingConceptComponent;
  let fixture: ComponentFixture<ProgrammingConceptComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProgrammingConceptComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProgrammingConceptComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
