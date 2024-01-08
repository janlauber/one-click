/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db);
  const collection = dao.findCollectionByNameOrId("w9y791qv1ymfu3m");

  return dao.deleteCollection(collection);
}, (db) => {
  const collection = new Collection({
    "id": "w9y791qv1ymfu3m",
    "created": "2024-01-08 11:03:03.128Z",
    "updated": "2024-01-08 11:09:38.816Z",
    "name": "settings",
    "type": "base",
    "system": false,
    "schema": [
      {
        "system": false,
        "id": "p3pryxsq",
        "name": "project",
        "type": "relation",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "7kff2zw80a7rmbu",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "wmfr90o2",
        "name": "blueprint",
        "type": "relation",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {
          "collectionId": "vs5gr49hpah1g9q",
          "cascadeDelete": false,
          "minSelect": null,
          "maxSelect": 1,
          "displayFields": null
        }
      },
      {
        "system": false,
        "id": "jodwbicb",
        "name": "overview",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "rde8xuub",
        "name": "rollouts",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "6sjgcuyl",
        "name": "image",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "zyaygdld",
        "name": "scale",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "x9r4zkcu",
        "name": "network",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "gfr3adhu",
        "name": "volumes",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "a00cz8lc",
        "name": "instances",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "90d6kc7n",
        "name": "envs",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      },
      {
        "system": false,
        "id": "axybtda3",
        "name": "settings",
        "type": "bool",
        "required": false,
        "presentable": false,
        "unique": false,
        "options": {}
      }
    ],
    "indexes": [],
    "listRule": "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id ?= blueprint.users.id || @request.auth.id = project.user.id)",
    "viewRule": "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id ?= blueprint.users.id || @request.auth.id = project.user.id)",
    "createRule": "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id = project.user.id)",
    "updateRule": "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id = project.user.id)",
    "deleteRule": "@request.auth.id != \"\" && (@request.auth.id = blueprint.owner.id || @request.auth.id = project.user.id)",
    "options": {}
  });

  return Dao(db).saveCollection(collection);
})
