import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ModuleConceptNodeComponent } from './module-concept-node.component';

describe('ModuleConceptNodeComponent', () => {
  let component: ModuleConceptNodeComponent;
  let fixture: ComponentFixture<ModuleConceptNodeComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ModuleConceptNodeComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ModuleConceptNodeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
