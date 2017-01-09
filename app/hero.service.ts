import { Injectable }     from '@angular/core';
import { Http, Response } from '@angular/http';
import { Hero }           from './hero';
import { Observable }     from 'rxjs/Observable';
import { Headers, RequestOptions } from '@angular/http';
@Injectable()
export class HeroService {
    private headers = new Headers({ 'Content-Type': 'application/json' });
  private heroesUrl = 'myhandle/heroes';  // URL to web API
  constructor (private http: Http) {}
  getHeroes (): Promise<Hero[]> {
    return this.http.get(this.heroesUrl)
                    .toPromise()
                    .then(this.extractData)
                    .catch(this.handleError);
  }
  private extractData(res: Response) {
    let body = res.json();
    return body || { };
  }
  private handleError (error: Response | any) {
    // In a real world app, we might use a remote logging infrastructure
    let errMsg: string;
    if (error instanceof Response) {
      const body = error.json() || '';
      const err = body.error || JSON.stringify(body);
      errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
    } else {
      errMsg = error.message ? error.message : error.toString();
    }
    console.error(errMsg);
    return Observable.throw(errMsg);
  }

  getHero(id: number): Promise<Hero> {
  return this.getHeroes()
        .then(heroes => heroes.find(hero => hero.id === id));
}


update(hero: Hero): Promise<Hero[]> {
  const url = `${this.heroesUrl}/${hero.id}`;
  return this.http
    .put(url, JSON.stringify(hero), {headers: this.headers})
    .toPromise()
    .then(() => hero)
    .catch(this.handleError);
}
delete(id: number): Promise<void> {
  const url = `${this.heroesUrl}/${id}`;
  return this.http.delete(url, {headers: this.headers})
    .toPromise()
    .then(() => null)
    .catch(this.handleError);
}
create(name: string): Promise<Hero> {
  return this.http
    .post(this.heroesUrl, JSON.stringify({name: name}), {headers: this.headers})
    .toPromise()
    .then(res => res.json())
    .catch(this.handleError);
}

}
