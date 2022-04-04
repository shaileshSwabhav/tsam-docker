import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BlogVerificationComponent } from './blog-verification.component';

describe('BlogVerificationComponent', () => {
  let component: BlogVerificationComponent;
  let fixture: ComponentFixture<BlogVerificationComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BlogVerificationComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BlogVerificationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
