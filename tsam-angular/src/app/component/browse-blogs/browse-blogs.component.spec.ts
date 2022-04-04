import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BrowseBlogsComponent } from './browse-blogs.component';

describe('BrowseBlogsComponent', () => {
  let component: BrowseBlogsComponent;
  let fixture: ComponentFixture<BrowseBlogsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BrowseBlogsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BrowseBlogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
