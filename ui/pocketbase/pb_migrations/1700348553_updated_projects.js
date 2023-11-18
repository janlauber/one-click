/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "c8q6kq0e",
    "name": "rollouts",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "22k6ts6gvnp46mc",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": null,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // update
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "c8q6kq0e",
    "name": "deployments",
    "type": "relation",
    "required": false,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "22k6ts6gvnp46mc",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": null,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
})
