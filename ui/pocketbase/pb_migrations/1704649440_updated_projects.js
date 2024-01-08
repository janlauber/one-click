/// <reference path="../pb_data/types.d.ts" />
migrate((db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // remove
  collection.schema.removeField("a504s0mr")

  return dao.saveCollection(collection)
}, (db) => {
  const dao = new Dao(db)
  const collection = dao.findCollectionByNameOrId("7kff2zw80a7rmbu")

  // add
  collection.schema.addField(new SchemaField({
    "system": false,
    "id": "a504s0mr",
    "name": "framework",
    "type": "relation",
    "required": true,
    "presentable": false,
    "unique": false,
    "options": {
      "collectionId": "nwgspeyl1n5v30f",
      "cascadeDelete": false,
      "minSelect": null,
      "maxSelect": 1,
      "displayFields": null
    }
  }))

  return dao.saveCollection(collection)
})
